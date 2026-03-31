package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jinwoll/minixia-cli/internal/api"
	"github.com/jinwoll/minixia-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	mqttBroker   string
	mqttUsername  string
	mqttPassword string
	mqttExec     string
)

var mqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "通过 MQTT 订阅指令",
	Long: `连接 MQTT Broker 并订阅指令 topic，实时接收用户下发的命令。

示例：
  minixia mqtt --broker mqtt://broker.minixia.app:1883
  minixia mqtt --broker mqtt://broker.minixia.app:1883 --exec './handle.sh "$CONTENT"'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.Apikey == "" {
			return fmt.Errorf("缺少 API Key，请运行 minixia init 或通过 --apikey 指定")
		}
		if mqttBroker == "" {
			return fmt.Errorf("请通过 --broker 指定 MQTT Broker 地址")
		}

		// 订阅的 topic 格式：cmd/{apikey}/{role}
		topic := fmt.Sprintf("cmd/%s/%s", resolvedCfg.Apikey, resolvedCfg.Role)
		output.PrintInfo(fmt.Sprintf("正在连接 %s …", mqttBroker))

		opts := mqtt.NewClientOptions().
			AddBroker(mqttBroker).
			SetClientID(fmt.Sprintf("minixia-cli-%d", time.Now().UnixMilli())).
			SetAutoReconnect(true).
			SetConnectRetry(true).
			SetConnectRetryInterval(5 * time.Second)

		if mqttUsername != "" {
			opts.SetUsername(mqttUsername)
		}
		if mqttPassword != "" {
			opts.SetPassword(mqttPassword)
		}

		// 连接断开时输出提示
		opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
			output.PrintWarning(fmt.Sprintf("MQTT 连接断开: %v，正在重连…", err))
		})
		opts.SetOnConnectHandler(func(c mqtt.Client) {
			output.PrintSuccess(fmt.Sprintf("已连接，订阅 topic: %s", topic))
			// 每次重连后重新订阅
			c.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
				handleMqttMessage(msg)
			})
		})

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return fmt.Errorf("MQTT 连接失败: %w", token.Error())
		}

		// 阻塞等待退出信号
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		fmt.Println()
		output.PrintInfo("正在断开 MQTT 连接…")
		client.Disconnect(1000)
		output.PrintInfo("已退出")
		return nil
	},
}

// handleMqttMessage 处理收到的 MQTT 指令消息
func handleMqttMessage(msg mqtt.Message) {
	var command api.Command
	if err := json.Unmarshal(msg.Payload(), &command); err != nil {
		output.PrintWarning(fmt.Sprintf("解析指令消息失败: %v", err))
		fmt.Println(string(msg.Payload()))
		return
	}

	fmt.Printf("[%s] %s (%s): %s\n",
		output.FormatTimestamp(command.Timestamp),
		command.ClientCmdID, command.Type, command.Content)

	if mqttExec != "" {
		executeCommand(mqttExec, command.Content, command.Type, command.ClientCmdID)
	}
}

func init() {
	mqttCmd.Flags().StringVarP(&mqttBroker, "broker", "b", "", "MQTT Broker 地址")
	mqttCmd.Flags().StringVar(&mqttUsername, "username", "", "MQTT 用户名")
	mqttCmd.Flags().StringVar(&mqttPassword, "password", "", "MQTT 密码")
	mqttCmd.Flags().StringVarP(&mqttExec, "exec", "e", "", "对每条指令执行的 shell 命令")
	rootCmd.AddCommand(mqttCmd)
}
