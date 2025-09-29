package features

import (
	"github.com/cucumber/godog"
)

func StepDefinitioninition1(arg1 *godog.Table) error {
	return godog.ErrPending
}

func StepDefinitioninition2(arg1 int) error {
	return godog.ErrPending
}

func StepDefinitioninition3(arg1 int) error {
	return godog.ErrPending
}

func StepDefinitioninition4(arg1 string) error {
	return godog.ErrPending
}

func gET(arg1 string) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^データベースに以下のユーザーが存在する:$`, StepDefinitioninition1)
	ctx.Step(`^レスポンスのステータスコードは(\d+)であるべき$`, StepDefinitioninition2)
	ctx.Step(`^レスポンスのボディには(\d+)人のユーザーが含まれているべき$`, StepDefinitioninition3)
	ctx.Step(`^ユーザー名 "([^"]*)" がレスポンスに含まれているべき$`, StepDefinitioninition4)
	ctx.Step(`^"([^"]*)" にGETリクエストを送信する$`, gET)
}
