/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : main.go
#   Created       : 2019-01-08 14:26:16
#   Last Modified : 2019-01-08 14:26:16
#   Describe      :
#
# ====================================================*/
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
	"uuabc.com/sendmsg/api"
	"uuabc.com/sendmsg/api/version"
	"uuabc.com/sendmsg/config"
	"uuabc.com/sendmsg/pkg/log"
	"uuabc.com/sendmsg/pkg/opentracing"

	"github.com/spf13/cobra"
	"uuabc.com/sendmsg/pkg/cmd"
)

var (
	opts = &Options{
		host:       "0.0.0.0",
		port:       8990,
		configPath: "/app/sendmsg/conf/config.yaml",
		logPath:    "/app/sendmsg/log/log.log",
		logLevel:   "info",
	}

	rootCmd = &cobra.Command{
		Use:          "receive-server",
		Short:        "The producer of the message service.",
		SilenceUsage: true,
	}

	startCmd = &cobra.Command{
		Use: "start",
		Long: `Start the service to receive the parameters from the user 
		and send the parameters to mq for consumption by the consumer`,
		RunE: start,
	}
)

const (
	defaultTimeout = time.Second * 10
)

type Options struct {
	host        string
	port        int
	configPath  string
	logPath     string
	logLevel    string
	addrJaeger  string
	addrMonitor string
}

func init() {
	startCmd.PersistentFlags().StringVarP(&opts.host, "host", "s", opts.host, "host for service startup")
	startCmd.PersistentFlags().IntVarP(&opts.port, "port", "p", opts.port, "port for service startup")
	startCmd.PersistentFlags().StringVarP(&opts.configPath, "config-path", "f", opts.configPath, "the path of the config file")
	startCmd.PersistentFlags().StringVar(&opts.logPath, "log-path", opts.logPath, "the location of the log file output")
	startCmd.PersistentFlags().StringVar(&opts.logLevel, "log-level", opts.logLevel, "log file output level")
	startCmd.PersistentFlags().StringVar(&opts.addrJaeger, "addr-jaeger", opts.addrJaeger, "the address of jaeger")
	startCmd.PersistentFlags().StringVar(&opts.addrMonitor, "addr-monitor", opts.addrMonitor, "the address of monitor(prometheus)")

	cmd.AddFlags(rootCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(version.Command())
}

func start(_ *cobra.Command, _ []string) error {
	stopC := cmd.GratefulQuit()
	var err error
	printFlags()

	// init log
	log.Init(opts.logPath, opts.logLevel)

	if err = config.Init(opts.configPath); err != nil {
		return err
	}

	// init opentracing
	if err = opentracing.New(opentracing.InitConfig(opts.addrJaeger)).Setup(); err != nil {
		return err
	}

	r := mux.NewRouter()

	if err := api.Init(r, opts.addrMonitor); err != nil {
		return err
	}

	svr := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", opts.host, opts.port),
		Handler:      r,
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
		IdleTimeout:  defaultTimeout,
	}

	// grateful quit
	go func() {
		<-stopC
		logrus.Info("stopping server now")
		if err := svr.Close(); err != nil {
			logrus.Errorf("Server Close:", err)
		}
	}()
	// start server
	if err = svr.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logrus.Info("Server closed under request\n")
			return nil
		} else {
			logrus.Infof("Server closed unexpect, %s\n", err.Error())
		}
	}
	return err
}

func printFlags() {
	logrus.WithField("Host", opts.host).Info()
	logrus.WithField("Post", opts.port).Info()
	logrus.WithField("Config-Path", opts.configPath).Info()
	logrus.WithField("Log-Path", opts.logPath).Info()
	logrus.WithField("Log-Level", opts.logLevel).Info()
	logrus.WithField("Addr-Jaeger", opts.addrJaeger).Info()
	logrus.WithField("Addr-Monitor", opts.addrMonitor).Info()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("serve failed ,error: %v\n", err)
		os.Exit(-1)
	}
}