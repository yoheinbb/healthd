package e2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Ret struct {
	Result string
}

func newRet(body []byte) (Ret, error) {
	var ret Ret
	if err := json.Unmarshal(body, &ret); err != nil {
		return Ret{}, err
	}
	return ret, nil
}

func upHealthdSuccess(ctx context.Context) {
	_, err := exec.CommandContext(ctx,
		"make", "up-success").CombinedOutput()
	Expect(err).ShouldNot(HaveOccurred())
}
func upHealthdFail(ctx context.Context) {
	_, err := exec.CommandContext(ctx,
		"make", "up-fail").CombinedOutput()
	Expect(err).ShouldNot(HaveOccurred())
}
func upHealthdTimeout(ctx context.Context) {
	_, err := exec.CommandContext(ctx,
		"make", "up-timeout").CombinedOutput()
	Expect(err).ShouldNot(HaveOccurred())
}
func downHealthdSuccess(ctx context.Context) {
	_, _ = exec.CommandContext(ctx,
		"make", "down-success").CombinedOutput()
}
func downHealthdFail(ctx context.Context) {
	_, _ = exec.CommandContext(ctx,
		"make", "down-fail").CombinedOutput()
}
func downHealthdTimeout(ctx context.Context) {
	_, _ = exec.CommandContext(ctx,
		"make", "down-timeout").CombinedOutput()
}

func getHealthcheck(ctx context.Context, requestURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{Timeout: 1 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	Expect(err).ShouldNot(HaveOccurred())
	if err != nil {
		return nil, err
	}
	return body, nil
}

var _ = Describe("healthdのe2eテスト", Serial, func() {
	Context("正常系", func() {
		requestURL := "http://localhost:8080/healthcheck"

		It("scriptの終了コードが0の場合、成功になる", func() {
			ctx := context.Background()
			upHealthdSuccess(ctx)
			time.Sleep(1 * time.Second)

			body, err := getHealthcheck(ctx, requestURL)
			Expect(err).ShouldNot(HaveOccurred())

			ret, err := newRet(body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ret.Result).Should(Equal("SUCCESS"))

		})
		It("メンテナンステキスト配置した場合、失敗になる", func() {
			ctx := context.Background()
			_, err := exec.CommandContext(ctx,
				"make", "enter-maintenance").CombinedOutput()
			time.Sleep(1 * time.Second)

			Expect(err).ShouldNot(HaveOccurred())
			body, err := getHealthcheck(ctx, requestURL)
			Expect(err).ShouldNot(HaveOccurred())

			ret, err := newRet(body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ret.Result).Should(Equal("FAILED"))
		})
		It("メンテナンステキスト削除した場合、成功になる", func() {
			ctx := context.Background()
			_, err := exec.CommandContext(ctx,
				"make", "exit-maintenance").CombinedOutput()
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(1 * time.Second)

			body, err := getHealthcheck(ctx, requestURL)
			Expect(err).ShouldNot(HaveOccurred())

			ret, err := newRet(body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ret.Result).Should(Equal("SUCCESS"))
		})
		It("SIGHUPでログローテートされる", func() {
			ctx := context.Background()
			_, err := exec.CommandContext(ctx,
				"make", "sighup-healthd").CombinedOutput()
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(1 * time.Second)

			byteAll, err := exec.CommandContext(ctx,
				"make", "-s", "count-healthd-log").CombinedOutput()
			Expect(err).ShouldNot(HaveOccurred())
			time.Sleep(1 * time.Second)

			Expect(string(byteAll)).Should(Equal("2\n"))
		})

	})
	It("scriptの終了コードが0以外の場合、失敗になる", func() {
		requestURL := "http://localhost:8081/healthcheck"
		ctx := context.Background()
		upHealthdFail(ctx)
		time.Sleep(1 * time.Second)

		body, err := getHealthcheck(ctx, requestURL)
		Expect(err).ShouldNot(HaveOccurred())

		ret, err := newRet(body)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ret.Result).Should(Equal("FAILED"))
	})
	It("scriptがタイムアウトした場合、失敗になる", func() {
		requestURL := "http://localhost:8082/healthcheck"
		ctx := context.Background()
		upHealthdTimeout(ctx)
		time.Sleep(1 * time.Second)

		body, err := getHealthcheck(ctx, requestURL)
		Expect(err).ShouldNot(HaveOccurred())

		ret, err := newRet(body)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ret.Result).Should(Equal("FAILED"))

		time.Sleep(2 * time.Second)
		byteAll, err := exec.CommandContext(ctx,
			"make", "-s", "cat-healthd-log").CombinedOutput()
		Expect(err).ShouldNot(HaveOccurred())

		Expect(string(byteAll)).Should(ContainSubstring("timeout!! kill process"))
	})
})

var _ = BeforeSuite(func() {
	downHealthdFail(context.Background())
	downHealthdSuccess(context.Background())
	downHealthdTimeout(context.Background())
})

var _ = AfterSuite(func() {
	downHealthdFail(context.Background())
	downHealthdSuccess(context.Background())
	downHealthdTimeout(context.Background())
})
