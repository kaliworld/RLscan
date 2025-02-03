package Q_learning

import (
	"RLscan/pkg/RL/Progress"
	"RLscan/pkg/RL/RLScan"
	"RLscan/pkg/RL/utlis"
	"RLscan/pkg/common"
	"bufio"
	_ "embed"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
Q-learning算法基于数组所以未构建神经神经网络，没有记忆回廊和经验回放的能力。每次使用都要用贝尔曼方程产出Q值，效果是原Fscan速度的1.3-1.7左右波动，端口覆盖率为87.6%
*/

// Environment 模拟代理的环境
type Environment struct {
	numbers []int
	score   int
}

// NewEnvironment 创建一个具有随机数字和零分数的新环境
func NewEnvironment() *Environment {
	return &Environment{
		numbers: make([]int, 0),
		score:   0,
	}
}

type QAgent struct {
	QTable    map[int]map[int]float64 // QTable[state][action] = value
	Alpha     float64                 // 学习率
	Gamma     float64                 // 折现因子
	Epsilon   float64                 // 探索率
	Actions   []int                   // 可能的动作集
	StateSize int                     // 状态的数量
}

// NewQAgent 创建一个新的Q-Learning智能体
func NewQAgent(stateSize int, actions []int, alpha float64, gamma float64, epsilon float64) *QAgent {
	qTable := make(map[int]map[int]float64)
	for i := 0; i < stateSize; i++ {
		qValues := make(map[int]float64)
		for _, action := range actions {
			qValues[action] = 0.0 // 初始值为0
		}
		qTable[i] = qValues
	}

	return &QAgent{
		QTable:    qTable,
		Alpha:     alpha,
		Gamma:     gamma,
		Epsilon:   epsilon,
		Actions:   actions,
		StateSize: stateSize,
	}
}

//go:embed q_vul10000.txt
var Qvul string

// SelectAction 基于当前状态选择最优动作或者随机动作
func (agent *QAgent) SelectAction(state int) int {
	// 从知识中获取动作空间
	reader := strings.NewReader(Qvul)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		number, _ := strconv.Atoi(line)
		agent.Actions = append(agent.Actions, number)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "读取错误:", err)
	}

	// Epsilon-greedy策略
	if rand.Float64() < agent.Epsilon {
		return agent.Actions[rand.Intn(len(agent.Actions))]
	}

	// 选择具有最高Q值的动作
	bestAction := agent.Actions[0]
	maxValue := -math.MaxFloat64
	for a, value := range agent.QTable[state] {
		if value > maxValue {
			maxValue = value
			bestAction = a
		}
	}

	return bestAction
}

// learn 更新Q表
func (agent *QAgent) learn(state int, action int, reward int, nextState int) {
	// 确保给定状态存在于Q表中
	if _, ok := agent.QTable[state]; !ok {
		agent.QTable[state] = make(map[int]float64)
	}

	// 确保给定下一个状态也存在于Q表中
	if _, ok := agent.QTable[nextState]; !ok {
		agent.QTable[nextState] = make(map[int]float64)
	}

	// 继续之前的操作...
	oldValue := agent.QTable[state][action]
	nextMax := -math.MaxFloat64
	for _, val := range agent.QTable[nextState] {
		if val > nextMax {
			nextMax = val
		}
	}
	// 贝尔曼方程
	// 动作状态价值的转换关系
	newValue := oldValue + agent.Alpha*(float64(reward)+agent.Gamma*nextMax-oldValue)
	agent.QTable[state][action] = newValue
}

// Step 评估智能体的猜测并更新分数
func (env *Environment) Step(guess int) (int, int, bool) {
	var reward int // 分数
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(env.numbers))
	randomNumber := env.numbers[randomIndex]
	if guess == randomNumber {
		reward = 6 // 猜对得分
	} else {
		reward = -1 // 猜错得分
	}
	env.score += reward
	done := env.score >= 80 || env.score <= -30
	nextState := rand.Intn(len(env.numbers))
	return reward, nextState, done
}

// Reset 使用新的随机数初始化环境并重置分数
func (env *Environment) Reset() int {
	//rand.Seed(time.Now().UnixNano())
	//env.numbers = rand.Perm(65535)[:10] // 0分，范围太大且没有关联性

	file, _ := os.Open("pkg/RL/Q-learning/ProtTest.txt") // 一类设备中的端口
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		numberStrings := strings.Split(line, ",")
		for _, numberString := range numberStrings {
			number, _ := strconv.Atoi(numberString)
			env.numbers = append(env.numbers, number)
		}
	}
	env.score = 0
	return env.numbers[rand.Intn(10)]
}

// RandomAgent 模拟进行随机猜测的代理
type RandomAgent struct{}

// PredictOpenPorts 根据Q表预测可能打开的端口
func (agent *QAgent) PredictOpenPorts() []int {
	var openPorts []int

	// 对于每个状态，选择具有最高Q值的端口
	for state := range agent.QTable {
		maxValue := -math.MaxFloat64
		var bestAction int
		for action, value := range agent.QTable[state] {
			if value > maxValue {
				maxValue = value
				bestAction = action
			}
		}
		// 将最佳动作（预测为打开的端口）添加到结果中
		openPorts = append(openPorts, bestAction)
	}
	return openPorts
}

// Train1 Train 针对新的QAgent类更新Train函数
func Train1() {
	env := NewEnvironment()
	stateSize := 10
	actions := env.numbers

	alpha := 0.9   // 学习率 ，这个值范围是0-1，值越靠近0，越考虑眼前利益，值越靠近1越考虑长远利用
	gamma := 0.6   // 折扣因子
	epsilon := 0.7 // 探索率 , 这个值范围是0-1, 值越靠近0，越会利用已学习的知识，越靠近1越会探索
	agent := NewQAgent(stateSize, actions, alpha, gamma, epsilon)

	episodes := 2000 // 模拟的游戏次数

	for episode := 0; episode < episodes; episode++ {
		state := env.Reset()
		done := false     // 表示游戏是否结束
		totalReward := 20 // 总奖励

		for !done {
			action := agent.SelectAction(state)
			reward, nextState, finished := env.Step(action)

			agent.learn(state, action, reward, nextState)

			state = nextState
			totalReward += reward
			done = finished

			if done {
				break
			}
		}

		fmt.Printf("Episode %d: Total Reward: %d\n", episode, totalReward)
		// 如果达到目标分数，则结束训练
		if totalReward >= 80 {
			fmt.Println("Congratulations! The agent has learned enough to pass.")
			break
		}

		// 减少探索率
		if agent.Epsilon > 0.01 {
			agent.Epsilon *= 0.995
		}
	}

	//fmt.Println("Predicted open ports:", utlis.RemoveZeros(utlis.UniqueInts(agent.PredictOpenPorts())))
}

// Run 对接真实环境
func Run(Hosts []string, Timeout int64) []string {
	var reward int
	var nextState int
	var finished bool
	var ret string

	var actions []int
	stateSize := 10
	alpha := 0.9
	gamma := 0.6
	epsilon := 0.7
	state := 1
	agent := NewQAgent(stateSize, actions, alpha, gamma, epsilon)
	episodes := 1000

	bar := Progress.NewProgressBar(episodes, 50) // 总数为100，进度条长50个字符
	for episode := 0; episode < episodes; episode++ {
		done := false     // 表示游戏是否结束
		totalReward := 20 // 总奖励

		for !done {
			action := agent.SelectAction(state)

			Ports := strconv.Itoa(action)
			_, reward, nextState, finished, ret = RLScan.PortScan(Hosts, Ports, Timeout, totalReward)

			agent.learn(state, action, reward, nextState)

			state = nextState
			totalReward += reward
			done = finished

			if done {
				break
			}
		}

		bar.Current = episode
		bar.Show()
		//fmt.Printf("Episode %d: Total Reward: %d\n", episode, totalReward)
		// 如果达到目标分数，则结束训练
		if totalReward >= 80 {
			fmt.Println("Congratulations! The agent has learned enough to pass.")
			break
		}

		// 减少探索率
		if agent.Epsilon > 0.01 {
			agent.Epsilon *= 0.995
		}
	}
	fmt.Println(ret)
	common.LogSuccess(ret)
	return utlis.LinesFromString(ret)
}
