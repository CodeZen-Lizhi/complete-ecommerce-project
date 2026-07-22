package ai

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var exerciseDirectoryPattern = regexp.MustCompile(`^\d{2}_[a-z0-9_]+$`)
var todoPattern = regexp.MustCompile(`TODO ([0-9]+)`)

var requiredExerciseSymbols = map[string][]string{
	"phase2/13_customer_support_service":       {"chatProvider", "historyStore", "customerSupportService", "newCustomerSupportRouter", "defaultBusinessKnowledge", "newOpenAIProvider", "newRedisHistoryStore", "buildClassificationResponseFormat"},
	"phase3/01_embedding_similarity":           {"embedder", "parseEmbeddingInput", "loadEmbeddingConfig", "cosineSimilarity", "rankTopK"},
	"phase3/02_document_loading_metadata":      {"documentLoader", "pdfLoaderConfig", "newPDFPageLoader", "validateDocumentPath", "normalizeDocument"},
	"phase3/03_chunking_strategies":            {"chunker", "buildParentChildChunks"},
	"phase3/04_vector_index_and_retrieval":     {"indexer", "retriever", "qdrantConfig", "newQdrantHTTPClient", "validateIndexBatch"},
	"phase3/05_hnsw_tuning":                    {"hnswIndex", "evaluateHNSW"},
	"phase3/06_parent_child_retrieval":         {"childRetriever", "parentStore", "expandParents"},
	"phase3/07_hybrid_search_rrf":              {"denseRetriever", "sparseRetriever", "hybridBackendConfig", "newHybridRetrievers", "reciprocalRankFusion"},
	"phase3/08_query_rewrite_and_filter":       {"queryRewriter", "filteredRetriever", "validateQueryFilter"},
	"phase3/09_rerank":                         {"reranker", "rerankerConfig", "newRemoteReranker", "validateRerankResults"},
	"phase3/10_context_budget_and_citations":   {"tokenCounter", "selectContext"},
	"phase3/11_index_lifecycle":                {"documentStateStore", "vectorIndex", "newMySQLStateDB", "runImport"},
	"phase3/12_rag_evaluation":                 {"recallAtK", "reciprocalRank", "summarizeRAGEvaluation"},
	"phase4/01_local_readonly_tool":            {"readonlyTool", "toolExecutor", "newEinoToolsNode", "validateToolInfo"},
	"phase4/02_tool_argument_validation":       {"decodeInventoryQuery", "validateInventoryQuery", "validationError"},
	"phase4/03_context_identity_authorization": {"withIdentity", "identityFromContext", "authorizeResource"},
	"phase4/04_ecommerce_query_tools":          {"productQueryService", "orderQueryService", "queryOrderTool"},
	"phase4/05_multi_tool_orchestration":       {"toolRegistry", "modelTurn", "bindRealTools", "runToolLoop"},
	"phase4/06_tool_resilience":                {"classifyToolFailure", "invokeWithRetry"},
	"phase4/07_safe_write_tool":                {"confirmationStore", "idempotencyStore", "executeWriteTool"},
	"phase5/01_chain_workflow":                 {"chainNode", "buildChain"},
	"phase5/02_graph_routing":                  {"classifyRoute", "buildRoutingGraph"},
	"phase5/03_react_agent":                    {"reactModel", "reactToolsNode", "newEinoReActAgent", "runReAct"},
	"phase5/04_agent_state_and_budget":         {"stateStore", "checkBudget", "advanceState"},
	"phase5/05_loop_detection":                 {"normalizeArguments", "recordCall"},
	"phase5/06_interrupt_checkpoint_resume":    {"checkpointStore", "compileInterruptibleGraph", "resumeFromCheckpoint"},
	"phase5/07_async_agent_task":               {"taskRepository", "runWorker"},
	"phase5/08_task_recovery_idempotency":      {"recoveryStore", "idempotencyRecordStore", "recoverTasks"},
	"phase6/01_golden_dataset":                 {"loadGoldenDataset", "validateGoldenCase"},
	"phase6/02_retrieval_metrics":              {"recallAtK", "reciprocalRank", "summarizeRetrieval"},
	"phase6/03_answer_quality_checks":          {"judgeResult", "newJudgeModel", "validateCitations", "evaluateFacts"},
	"phase6/04_tool_and_agent_metrics":         {"evaluateExecution", "executionMetrics"},
	"phase6/05_latency_token_cost_metrics":     {"percentileDuration", "summarizeMeasurements"},
	"phase6/06_regression_comparison":          {"compareSnapshots", "regressionReport"},
	"phase6/07_eval_ci_gate":                   {"evaluateGate", "exitCodeForGate"},
	"phase7/01_secret_and_log_redaction":       {"loadSecretConfig", "redactFields", "containsSecret"},
	"phase7/02_prompt_injection_defense":       {"evaluateInjection", "authorizeToolCall"},
	"phase7/03_rate_limit_and_concurrency":     {"rateLimiter", "concurrencyLimiter", "runLimited"},
	"phase7/04_end_to_end_tracing":             {"tracer", "span", "newOpenTelemetryProvider", "tracePipeline"},
	"phase7/05_metrics_and_alerts":             {"metricRecorder", "validateMetricLabels", "evaluateAlert"},
	"phase7/06_docker_compose_stack":           {"validateComposeFile", "composeCommands"},
	"phase7/07_health_gray_release_rollback":   {"releaseControlConfig", "newReleaseControlClient", "validateHealth", "evaluateRollout"},
	"phase7/08_dependency_failure_fallback":    {"decideFallback", "recoveryController"},
}

var requiredExerciseConstants = map[string][]string{
	"phase3/01_embedding_similarity":         {"baseURL", "apiKey", "modelName"},
	"phase3/04_vector_index_and_retrieval":   {"embeddingBaseURL", "embeddingAPIKey", "embeddingModelName", "qdrantBaseURL", "qdrantAPIKey", "qdrantCollection"},
	"phase3/05_hnsw_tuning":                  {"qdrantBaseURL", "qdrantAPIKey", "qdrantCollection"},
	"phase3/06_parent_child_retrieval":       {"qdrantBaseURL", "qdrantAPIKey", "qdrantCollection"},
	"phase3/07_hybrid_search_rrf":            {"embeddingBaseURL", "embeddingAPIKey", "embeddingModelName", "qdrantBaseURL", "qdrantAPIKey", "qdrantCollection", "bm25BaseURL", "bm25APIKey"},
	"phase3/08_query_rewrite_and_filter":     {"modelBaseURL", "modelAPIKey", "modelName", "qdrantBaseURL", "qdrantAPIKey", "qdrantCollection"},
	"phase3/09_rerank":                       {"rerankerBaseURL", "rerankerAPIKey", "rerankerModelName"},
	"phase3/11_index_lifecycle":              {"mysqlDSN", "qdrantBaseURL", "qdrantAPIKey", "qdrantCollection"},
	"phase4/05_multi_tool_orchestration":     {"baseURL", "apiKey", "modelName"},
	"phase5/01_chain_workflow":               {"baseURL", "apiKey", "modelName"},
	"phase5/03_react_agent":                  {"baseURL", "apiKey", "modelName"},
	"phase6/03_answer_quality_checks":        {"baseURL", "apiKey", "modelName"},
	"phase7/07_health_gray_release_rollback": {"releaseControlBaseURL", "releaseControlAPIKey"},
}

// TestExerciseCatalog 验证七个阶段的目录数量、连续编号和必需骨架文件。
func TestExerciseCatalog(t *testing.T) {
	t.Helper()

	expectedCounts := map[int]int{
		1: 7,
		2: 13,
		3: 12,
		4: 7,
		5: 8,
		6: 7,
		7: 8,
	}

	for phase := 1; phase <= 7; phase++ {
		phase := phase
		t.Run(fmt.Sprintf("phase%d", phase), func(t *testing.T) {
			validatePhase(t, phase, expectedCounts[phase])
		})
	}
}

// validatePhase 验证单个阶段的目录顺序，以及 README 和 Go 骨架是否齐全。
func validatePhase(t *testing.T, phase int, expectedCount int) {
	t.Helper()

	phaseDir := fmt.Sprintf("phase%d", phase)
	entries, err := os.ReadDir(phaseDir)
	if err != nil {
		t.Fatalf("读取 %s 失败: %v", phaseDir, err)
	}

	directories := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !exerciseDirectoryPattern.MatchString(entry.Name()) {
			t.Fatalf("%s 包含不符合 NN_slug 规范的目录 %q", phaseDir, entry.Name())
		}
		directories = append(directories, entry.Name())
	}
	sort.Strings(directories)

	if len(directories) != expectedCount {
		t.Fatalf("%s 目录数量为 %d，期望 %d: %v", phaseDir, len(directories), expectedCount, directories)
	}

	for index, directory := range directories {
		expectedPrefix := fmt.Sprintf("%02d_", index+1)
		if directory[:3] != expectedPrefix {
			t.Errorf("%s 第 %d 个目录为 %q，期望前缀 %q", phaseDir, index+1, directory, expectedPrefix)
		}

		readmePath := filepath.Join(phaseDir, directory, "README.md")
		assertRegularFile(t, readmePath)
		if phase >= 2 {
			mainPath := filepath.Join(phaseDir, directory, "main.go")
			assertRegularFile(t, mainPath)
			key := filepath.Join(phaseDir, directory)

			todoPath := mainPath
			if phase >= 3 {
				exercisePath := filepath.Join(phaseDir, directory, "exercise.go")
				assertRegularFile(t, exercisePath)
				todoPath = exercisePath
			}
			switch key {
			case "phase2/10_structured_json_output":
				testPath := filepath.Join(phaseDir, directory, "main_test.go")
				assertRegularFile(t, testPath)
				validateTODOFileSet(t, mainPath, testPath)
				validateNonStackedTODOPlacement(t, mainPath)
				validateNonStackedTODOPlacement(t, testPath)
				validateTODODocumentation(t, readmePath, mainPath, testPath)
			case "phase2/13_customer_support_service":
				corePaths := []string{
					filepath.Join(phaseDir, directory, "provider.go"),
					filepath.Join(phaseDir, directory, "classification.go"),
					filepath.Join(phaseDir, directory, "history.go"),
					filepath.Join(phaseDir, directory, "governance.go"),
					filepath.Join(phaseDir, directory, "service.go"),
				}
				for _, corePath := range corePaths {
					assertRegularFile(t, corePath)
					validateNonStackedTODOPlacement(t, corePath)
				}
				testPath := filepath.Join(phaseDir, directory, "main_test.go")
				assertRegularFile(t, testPath)
				validateTODOFileSet(t, corePaths...)
				validateTODODocumentation(t, readmePath, corePaths...)
			default:
				validateTODOSequence(t, todoPath)
			}

			if phase >= 3 || key == "phase2/04_few_shot_comparison" {
				validateInlineTODOPlacement(t, todoPath)
				validateTODODocumentation(t, readmePath, todoPath)
			}
			if symbols, exists := requiredExerciseSymbols[key]; exists {
				validateExerciseSymbols(t, key, symbols)
			}
			if constants, exists := requiredExerciseConstants[key]; exists {
				validateExerciseConstants(t, key, constants)
			}
		}
	}
}

// validateExerciseConstants 验证外部依赖配置提供可直接修改的顶部占位常量。
func validateExerciseConstants(t *testing.T, directory string, required []string) {
	t.Helper()

	packages, err := parser.ParseDir(token.NewFileSet(), directory, func(info os.FileInfo) bool {
		return filepath.Ext(info.Name()) == ".go" && !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Errorf("解析练习目录 %s 失败: %v", directory, err)
		return
	}

	found := make(map[string]struct{})
	usageCount := make(map[string]int)
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if identifier, ok := node.(*ast.Ident); ok {
					usageCount[identifier.Name]++
				}
				return true
			})
			for _, declaration := range file.Decls {
				group, ok := declaration.(*ast.GenDecl)
				if !ok || group.Tok != token.CONST {
					continue
				}
				for _, specification := range group.Specs {
					value, ok := specification.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, name := range value.Names {
						found[name.Name] = struct{}{}
					}
				}
			}
		}
	}

	for _, constant := range required {
		if _, exists := found[constant]; !exists {
			t.Errorf("练习目录 %s 缺少顶部配置常量 %q", directory, constant)
			continue
		}
		if usageCount[constant] < 2 {
			t.Errorf("练习目录 %s 的配置常量 %q 只声明但未接入配置数据流", directory, constant)
		}
	}
}

// validateTODOFileSet 验证一个练习跨多个 Go 文件的 TODO 编号唯一且连续。
func validateTODOFileSet(t *testing.T, paths ...string) {
	t.Helper()

	contents := make([]byte, 0)
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("读取 TODO 文件 %s 失败: %v", path, err)
			return
		}
		contents = append(contents, content...)
	}

	matches := todoPattern.FindAllSubmatch(contents, -1)
	seen := make(map[int]struct{}, len(matches))
	numbers := make([]int, 0, len(matches))
	for _, match := range matches {
		number, err := strconv.Atoi(string(match[1]))
		if err != nil {
			t.Errorf("解析 TODO 编号失败: %v", err)
			return
		}
		if _, exists := seen[number]; exists {
			t.Errorf("%v 重复使用 TODO %d", paths, number)
			return
		}
		seen[number] = struct{}{}
		numbers = append(numbers, number)
	}
	sort.Ints(numbers)
	for index, number := range numbers {
		if expected := index + 1; number != expected {
			t.Errorf("%v 的 TODO 编号为 %v，缺少 TODO %d", paths, numbers, expected)
			return
		}
	}
}

// validateTODODocumentation 验证 README 步骤与源码填写点编号一致。
func validateTODODocumentation(t *testing.T, readmePath string, codePaths ...string) {
	t.Helper()

	readme, err := os.ReadFile(readmePath)
	if err != nil {
		t.Errorf("读取练习文档 %s 失败: %v", readmePath, err)
		return
	}
	code := make([]byte, 0)
	for _, codePath := range codePaths {
		content, readErr := os.ReadFile(codePath)
		if readErr != nil {
			t.Errorf("读取练习代码 %s 失败: %v", codePath, readErr)
			return
		}
		code = append(code, content...)
	}

	readmeNumbers := uniqueTODONumbers(readme)
	codeNumbers := uniqueTODONumbers(code)
	if strings.Join(readmeNumbers, ",") != strings.Join(codeNumbers, ",") {
		t.Errorf("README 与源码 TODO 编号不一致：%s=%v，%v=%v", readmePath, readmeNumbers, codePaths, codeNumbers)
	}
}

// uniqueTODONumbers 按首次出现顺序返回去重后的 TODO 编号。
func uniqueTODONumbers(content []byte) []string {
	matches := todoPattern.FindAllSubmatch(content, -1)
	seen := make(map[string]struct{}, len(matches))
	numbers := make([]string, 0, len(matches))
	for _, match := range matches {
		number := string(match[1])
		if _, exists := seen[number]; exists {
			continue
		}
		seen[number] = struct{}{}
		numbers = append(numbers, number)
	}
	sort.Slice(numbers, func(left int, right int) bool {
		leftNumber, _ := strconv.Atoi(numbers[left])
		rightNumber, _ := strconv.Atoi(numbers[right])
		return leftNumber < rightNumber
	})
	return numbers
}

// validateInlineTODOPlacement 验证练习步骤分散在可填写位置，而不是连续堆在入口函数中。
func validateInlineTODOPlacement(t *testing.T, path string) {
	t.Helper()
	validateNonStackedTODOPlacement(t, path)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("读取 TODO 文件 %s 失败: %v", path, err)
		return
	}

	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, path, content, parser.ParseComments)
	if err != nil {
		t.Errorf("解析 TODO 文件 %s 失败: %v", path, err)
		return
	}
	for _, declaration := range parsedFile.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			continue
		}
		start := fileSet.Position(function.Body.Pos()).Offset
		end := fileSet.Position(function.Body.End()).Offset
		if matches := todoPattern.FindAll(content[start:end], -1); len(matches) > 1 {
			t.Errorf("%s 的函数 %s 包含 %d 个编号 TODO，请拆成独立填写函数", path, function.Name.Name, len(matches))
			return
		}
	}
}

// validateNonStackedTODOPlacement 验证单文件中的 TODO 不连续堆叠且编号不重复。
func validateNonStackedTODOPlacement(t *testing.T, path string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("读取 TODO 文件 %s 失败: %v", path, err)
		return
	}

	lines := strings.Split(string(content), "\n")
	previousWasTODO := false
	for lineNumber, line := range lines {
		isTODO := todoPattern.MatchString(line)
		if isTODO && previousWasTODO {
			t.Errorf("%s:%d 存在连续堆叠的编号 TODO，请把步骤移到对应代码位置", path, lineNumber+1)
			return
		}
		previousWasTODO = isTODO
	}

	seenNumbers := make(map[string]struct{})
	for _, match := range todoPattern.FindAllSubmatch(content, -1) {
		number := string(match[1])
		if _, exists := seenNumbers[number]; exists {
			t.Errorf("%s 重复使用 TODO %s，请让文档步骤与源码填写点一一对应", path, number)
			return
		}
		seenNumbers[number] = struct{}{}
	}
}

// validateExerciseSymbols 验证练习包包含题目专属接口、类型或函数签名。
func validateExerciseSymbols(t *testing.T, directory string, required []string) {
	t.Helper()

	packages, err := parser.ParseDir(token.NewFileSet(), directory, func(info os.FileInfo) bool {
		return filepath.Ext(info.Name()) == ".go" && !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Errorf("解析练习目录 %s 失败: %v", directory, err)
		return
	}

	found := make(map[string]struct{})
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, declaration := range file.Decls {
				switch value := declaration.(type) {
				case *ast.FuncDecl:
					found[value.Name.Name] = struct{}{}
				case *ast.GenDecl:
					for _, specification := range value.Specs {
						if typeSpec, ok := specification.(*ast.TypeSpec); ok {
							found[typeSpec.Name.Name] = struct{}{}
						}
					}
				}
			}
		}
	}

	for _, symbol := range required {
		if _, exists := found[symbol]; !exists {
			t.Errorf("练习目录 %s 缺少题目专属符号 %q", directory, symbol)
		}
	}
}

// validateTODOSequence 验证 Go 骨架中首次出现的 TODO 编号从 1 连续递增。
func validateTODOSequence(t *testing.T, path string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("读取 TODO 文件 %s 失败: %v", path, err)
		return
	}

	matches := todoPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		t.Errorf("Go 骨架 %s 没有编号 TODO", path)
		return
	}

	seen := make(map[int]struct{}, len(matches))
	numbers := make([]int, 0, len(matches))
	for _, match := range matches {
		number, err := strconv.Atoi(string(match[1]))
		if err != nil {
			t.Errorf("解析 %s 的 TODO 编号失败: %v", path, err)
			return
		}
		if _, exists := seen[number]; exists {
			continue
		}
		seen[number] = struct{}{}
		numbers = append(numbers, number)
	}

	for index, number := range numbers {
		expected := index + 1
		if number != expected {
			t.Errorf("%s 的首次 TODO 顺序为 %v，期望在位置 %d 出现 TODO %d", path, numbers, index+1, expected)
			return
		}
	}
}

// assertRegularFile 验证指定路径存在且是普通文件。
func assertRegularFile(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("必需文件 %s 不可用: %v", path, err)
		return
	}
	if !info.Mode().IsRegular() {
		t.Errorf("必需路径 %s 不是普通文件", path)
	}
}
