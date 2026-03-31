package cmd

import (
	"fmt"

	"github.com/jinwoll/minixia-cli/internal/api"
	"github.com/jinwoll/minixia-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	webhookURL    string
	webhookRemove bool
	webhookTest   bool
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "管理 Webhook 配置",
	Long: `设置、测试或移除 Webhook 回调地址。

示例：
  minixia webhook --url https://my-server.com/webhook
  minixia webhook --test
  minixia webhook --remove`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resolvedCfg.Apikey == "" {
			return fmt.Errorf("缺少 API Key，请运行 minixia init 或通过 --apikey 指定")
		}

		// 移除 Webhook
		if webhookRemove {
			if err := apiClient.RemoveWebhook(resolvedCfg.Apikey); err != nil {
				output.HandleError(err)
				return nil
			}
			output.PrintSuccess("Webhook 已移除")
			return nil
		}

		// 测试 Webhook
		if webhookTest {
			if err := apiClient.TestWebhook(resolvedCfg.Apikey); err != nil {
				output.HandleError(err)
				return nil
			}
			output.PrintSuccess("Webhook 测试回调已发送，请检查你的服务器是否收到。")
			return nil
		}

		// 设置 Webhook
		if webhookURL != "" {
			req := &api.WebhookSetRequest{
				Apikey: resolvedCfg.Apikey,
				Role:   resolvedCfg.Role,
				URL:    webhookURL,
			}
			if err := apiClient.SetWebhook(req); err != nil {
				output.HandleError(err)
				return nil
			}
			output.PrintSuccess(fmt.Sprintf("Webhook 已设置: %s", webhookURL))
			return nil
		}

		// 无 flag 时显示当前配置
		cfg, err := apiClient.GetWebhook(resolvedCfg.Apikey)
		if err != nil {
			output.HandleError(err)
			return nil
		}
		if cfg == nil {
			output.PrintInfo("尚未配置 Webhook。使用 --url 设置。")
			return nil
		}
		output.PrintKeyValue([][]string{
			{"URL", cfg.URL},
			{"已验证", fmt.Sprintf("%v", cfg.IsVerified)},
			{"重试次数", fmt.Sprintf("%d", cfg.RetryCount)},
			{"创建时间", cfg.CreatedAt},
		})
		return nil
	},
}

func init() {
	webhookCmd.Flags().StringVarP(&webhookURL, "url", "u", "", "Webhook 回调地址（HTTPS）")
	webhookCmd.Flags().BoolVar(&webhookRemove, "remove", false, "移除已配置的 Webhook")
	webhookCmd.Flags().BoolVar(&webhookTest, "test", false, "发送测试回调")
	rootCmd.AddCommand(webhookCmd)
}
