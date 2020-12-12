package flags

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DockerAPIMinVersion is the minimum version of the docker api required to
// use watchtower
const DockerAPIMinVersion string = "1.25"
const DefaultInterval = int(time.Hour * 24 / time.Second)

// RegisterDockerFlags that are used directly by the docker api client
func RegisterDockerFlags(rootCmd *cobra.Command) {
	flags := rootCmd.PersistentFlags()
	flags.StringP("host", "H", "unix:///var/run/docker.sock", "daemon socket to connect to")
	flags.BoolP("tlsverify", "v", false, "use TLS and verify the remote")
	flags.StringP("api-version", "a", DockerAPIMinVersion, "api version to use by docker client")
}

// RegisterSystemFlags that are used by watchtower to modify the program flow
func RegisterSystemFlags(rootCmd *cobra.Command) {
	flags := rootCmd.PersistentFlags()
	flags.IntP(
		"interval",
		"i",
		DefaultInterval, // viper.GetInt("WATCHTOWER_POLL_INTERVAL"),
		"poll interval (in seconds)")

	flags.String(
		"schedule",
		"s",
		"",
		"The cron expression which defines when to update")
	//viper.GetString("WATCHTOWER_SCHEDULE"),

	flags.DurationP(
		"stop-timeout",
		"t",
		time.Second*10, //viper.GetDuration("WATCHTOWER_TIMEOUT"),
		"Timeout before a container is forcefully stopped")

	flags.BoolP(
		"no-pull",
		"",
		false, // viper.GetBool("WATCHTOWER_NO_PULL"),
		"Do not pull any new images")

	flags.Bool(
		"no-restart",
		false, // viper.GetBool("WATCHTOWER_NO_RESTART"),
		"Do not restart any containers")

	flags.Bool(
		"no-startup-message",
		false, // viper.GetBool("WATCHTOWER_NO_STARTUP_MESSAGE"),
		"Prevents watchtower from sending a startup message")

	flags.BoolP(
		"cleanup",
		"c",
		false, // viper.GetBool("WATCHTOWER_CLEANUP"),
		"Remove previously used images after updating")

	flags.BoolP(
		"remove-volumes",
		"",
		false, // viper.GetBool("WATCHTOWER_REMOVE_VOLUMES"),
		"Remove attached volumes before updating")

	flags.BoolP(
		"label-enable",
		"e",
		false, // viper.GetBool("WATCHTOWER_LABEL_ENABLE"),
		"Watch containers where the com.centurylinklabs.watchtower.enable label is true")

	flags.BoolP(
		"debug",
		"d",
		false, // viper.GetBool("WATCHTOWER_DEBUG"),
		"Enable debug mode with verbose logging")

	flags.Bool(
		"trace",
		false, // viper.GetBool("WATCHTOWER_TRACE"),
		"Enable trace mode with very verbose logging - caution, exposes credentials")

	flags.BoolP(
		"monitor-only",
		"m",
		false, // viper.GetBool("WATCHTOWER_MONITOR_ONLY"),
		"Will only monitor for new images, not update the containers")

	flags.BoolP(
		"run-once",
		"R",
		false, // viper.GetBool("WATCHTOWER_RUN_ONCE"),
		"Run once now and exit")

	flags.BoolP(
		"include-restarting",
		"",
		false, // viper.GetBool("WATCHTOWER_INCLUDE_RESTARTING"),
		"Will also include restarting containers")

	flags.BoolP(
		"include-stopped",
		"S",
		false, // viper.GetBool("WATCHTOWER_INCLUDE_STOPPED"),
		"Will also include created and exited containers")

	flags.Bool(
		"revive-stopped",
		false, // viper.GetBool("WATCHTOWER_REVIVE_STOPPED"),
		"Will also start stopped containers that were updated, if include-stopped is active")

	flags.Bool(
		"enable-lifecycle-hooks",
		false, // viper.GetBool("WATCHTOWER_LIFECYCLE_HOOKS"),
		"Enable the execution of commands triggered by pre- and post-update lifecycle hooks")

	flags.Bool(
		"rolling-restart",
		false, // viper.GetBool("WATCHTOWER_ROLLING_RESTART"),
		"Restart containers one at a time")

	flags.Bool(
		"http-api-update",
		false, // viper.GetBool("WATCHTOWER_HTTP_API_UPDATE"),
		"Runs Watchtower in HTTP API mode, so that image updates must to be triggered by a request")
	flags.Bool(
		"http-api-metrics",
		false, // viper.GetBool("WATCHTOWER_HTTP_API_METRICS"),
		"Runs Watchtower with the Prometheus metrics API enabled")

	flags.String(
		"http-api-token",
		"", // viper.GetString("WATCHTOWER_HTTP_API_TOKEN"),
		"Sets an authentication token to HTTP API requests.")
	flags.BoolP(
		"http-api-periodic-polls",
		"",
		viper.GetBool("WATCHTOWER_HTTP_API_PERIODIC_POLLS"),
		"Also run periodic updates (specified with --interval and --schedule) if HTTP API is enabled")
	// https://no-color.org/
	flags.BoolP(
		"no-color",
		"",
		false, // viper.IsSet("NO_COLOR"),
		"Disable ANSI color escape codes in log output")
	flags.String(
		"scope",
		"", // viper.GetString("WATCHTOWER_SCOPE"),
		"Defines a monitoring scope for the Watchtower instance.")
}

// RegisterNotificationFlags that are used by watchtower to send notifications
func RegisterNotificationFlags(rootCmd *cobra.Command) {
	flags := rootCmd.PersistentFlags()

	flags.StringSliceP(
		"notifications",
		"n",
		[]string{}, // viper.GetStringSlice("WATCHTOWER_NOTIFICATIONS"),
		" Notification types to send (valid: email, slack, msteams, gotify, shoutrrr)")

	flags.String(
		"notifications-level",
		"info", // viper.GetString("WATCHTOWER_NOTIFICATIONS_LEVEL"),
		"The log level used for sending notifications. Possible values: panic, fatal, error, warn, info or debug")

	flags.Int(
		"notifications-delay",
		0, // viper.GetInt("WATCHTOWER_NOTIFICATIONS_DELAY"),
		"Delay before sending notifications, expressed in seconds")

	flags.String(
		"notifications-hostname",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATIONS_HOSTNAME"),
		"Custom hostname for notification titles")

	flags.String(
		"notification-email-from",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_EMAIL_FROM"),
		"Address to send notification emails from")

	flags.String(
		"notification-email-to",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_EMAIL_TO"),
		"Address to send notification emails to")

	flags.Int(
		"notification-email-delay",
		0,
		//viper.GetInt("WATCHTOWER_NOTIFICATION_EMAIL_DELAY"),
		"Delay before sending notifications, expressed in seconds")

	flags.String(
		"notification-email-server",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_EMAIL_SERVER"),
		"SMTP server to send notification emails through")

	flags.Int(
		"notification-email-server-port",
		25,
		// viper.GetInt("WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PORT"),
		"SMTP server port to send notification emails through")

	flags.Bool(
		"notification-email-server-tls-skip-verify",
		false,
		// viper.GetBool("WATCHTOWER_NOTIFICATION_EMAIL_SERVER_TLS_SKIP_VERIFY"),
		`Controls whether watchtower verifies the SMTP server's certificate chain and host name.
Should only be used for testing.`)

	flags.String(
		"notification-email-server-user",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_EMAIL_SERVER_USER"),
		"SMTP server user for sending notifications")

	flags.String(
		"notification-email-server-password",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PASSWORD"),
		"SMTP server password for sending notifications")

	flags.String(
		"notification-email-subjecttag",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_EMAIL_SUBJECTTAG"),
		"Subject prefix tag for notifications via mail")

	flags.String(
		"notification-slack-hook-url",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_SLACK_HOOK_URL"),
		"The Slack Hook URL to send notifications to")

	flags.String(
		"notification-slack-identifier",
		"watchtower",
		// viper.GetString("WATCHTOWER_NOTIFICATION_SLACK_IDENTIFIER"),
		"A string which will be used to identify the messages coming from this watchtower instance")

	flags.String(
		"notification-slack-channel",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_SLACK_CHANNEL"),
		"A string which overrides the webhook's default channel. Example: #my-custom-channel")

	flags.String(
		"notification-slack-icon-emoji",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_SLACK_ICON_EMOJI"),
		"An emoji code string to use in place of the default icon")

	flags.String(
		"notification-slack-icon-url",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_SLACK_ICON_URL"),
		"An icon image URL string to use in place of the default icon")

	flags.String(
		"notification-msteams-hook",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_MSTEAMS_HOOK_URL"),
		"The MSTeams WebHook URL to send notifications to")

	flags.Bool(
		"notification-msteams-data",
		false,
		// viper.GetBool("WATCHTOWER_NOTIFICATION_MSTEAMS_USE_LOG_DATA"),
		"The MSTeams notifier will try to extract log entry fields as MSTeams message facts")

	flags.String(
		"notification-gotify-url",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_GOTIFY_URL"),
		"The Gotify URL to send notifications to")

	flags.String(
		"notification-gotify-token",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_GOTIFY_TOKEN"),
		"The Gotify Application required to query the Gotify API")

	flags.Bool(
		"notification-gotify-tls-skip-verify",
		false,
		// viper.GetBool("WATCHTOWER_NOTIFICATION_GOTIFY_TLS_SKIP_VERIFY"),
		`Controls whether watchtower verifies the Gotify server's certificate chain and host name.
Should only be used for testing.`)

	flags.String(
		"notification-template",
		"",
		// viper.GetString("WATCHTOWER_NOTIFICATION_TEMPLATE"),
		"The shoutrrr text/template for the messages")

	flags.StringArray(
		"notification-url",
		[]string{},
		// viper.GetStringSlice("WATCHTOWER_NOTIFICATION_URL"),
		"The shoutrrr URL to send notifications to")

	flags.Bool("notification-report",
		false,
		// viper.GetBool("WATCHTOWER_NOTIFICATION_REPORT"),
		"Use the session report as the notification template data")

	flags.String(
		"warn-on-head-failure",
		"auto",
		// viper.GetString("WATCHTOWER_WARN_ON_HEAD_FAILURE"),
		"When to warn about HEAD pull requests failing. Possible values: always, auto or never")

}

// SetDefaults provides default values for environment variables
func SetDefaults() {
	day := (time.Hour * 24).Seconds()
	viper.AutomaticEnv()
}

// BindViperFlags binds the cmd PFlags to the viper configuration
func BindViperFlags(cmd *cobra.Command) {
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Fatalf("failed to bind flags: %v", err)
	}
}

// EnvConfig translates the command-line options into environment variables
// that will initialize the api client
func EnvConfig() error {
	var err error

	host := viper.GetString("host")
	tls = viper.GetBool("tlsverify")
	version = viper.GetString("api-version")
	if err = setEnvOptStr("DOCKER_HOST", host); err != nil {
		return err
	}
	if err = setEnvOptBool("DOCKER_TLS_VERIFY", tls); err != nil {
		return err
	}
	if err = setEnvOptStr("DOCKER_API_VERSION", version); err != nil {
		return err
	}
	return nil
}

// ReadFlags reads common flags used in the main program flow of watchtower
func ReadFlags() (cleanup bool, noRestart bool, monitorOnly bool, timeout time.Duration) {

	cleanup = viper.GetBool("cleanup")
	noRestart = viper.GetBool("no-restart")
	monitorOnly = viper.GetBool("monitor-only")
	timeout = viper.GetDuration("stop-timeout")

	return
}

func setEnvOptStr(env string, opt string) error {
	if opt == "" || opt == os.Getenv(env) {
		return nil
	}
	err := os.Setenv(env, opt)
	if err != nil {
		return err
	}
	return nil
}

func setEnvOptBool(env string, opt bool) error {
	if opt {
		return setEnvOptStr(env, "1")
	}
	return nil
}

// GetSecretsFromFiles checks if passwords/tokens/webhooks have been passed as a file instead of plaintext.
// If so, the value of the flag will be replaced with the contents of the file.
func GetSecretsFromFiles() {
	secrets := []string{
		"notification-email-server-password",
		"notification-slack-hook-url",
		"notification-msteams-hook",
		"notification-gotify-token",
	}
	for _, secret := range secrets {
		getSecretFromFile(secret)
	}
}

// getSecretFromFile will check if the flag contains a reference to a file; if it does, replaces the value of the flag with the contents of the file.
func getSecretFromFile(secret string) {
	value := viper.GetString(secret)
	if value != "" && isFile(value) {
		file, err := ioutil.ReadFile(value)
		if err != nil {
			log.Fatal(err)
		}
		viper.Set(secret, strings.TrimSpace(string(file)))
	}
}

func isFile(s string) bool {
	_, err := os.Stat(s)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
