package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Done      bool      `json:"done"`
	Limit     time.Time `json:"limit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	IsDone    bool      `json:"is_done"`
	IsDeleted bool      `json:"is_deleted"`
}

const fileName = "todos.json"

func loadTodos() ([]Todo, error) {
	var todos []Todo
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return todos, nil
		}
		return nil, err
	}
	err = json.Unmarshal(data, &todos)
	return todos, err
}

func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func add(task string, reader *bufio.Reader) {
	fmt.Println("📅 タスクの期限を入力してください (例: 2025/01/01 12:00:00)")
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	limit, err := time.Parse("2006/01/02 15:04:05", input) // 数字の順がアメリカだと123456になっていいい感じらしい
	if err != nil {
		fmt.Println("⚠️ 期限の形式が正しくありません")
		return
	}

	todos, _ := loadTodos()
	id := 1
	if len(todos) > 0 {
		id = todos[len(todos)-1].ID + 1
	}
	todo := Todo{ID: id, Task: task, Done: false, Limit: limit, CreatedAt: time.Now(), UpdatedAt: time.Now(), DeletedAt: time.Time{}, IsDone: false, IsDeleted: false}
	todos = append(todos, todo)
	saveTodos(todos)
	fmt.Println("✅ Todo追加:", task)
}

func format(id int, status string, task string, limit time.Time) string {
	len_status := 12 - utf8.RuneCountInString(status)
	len_task := 24 - utf8.RuneCountInString(task)
	if utf8.RuneCountInString(status) == len(status) {
		len_status += len(status)
	}
	if utf8.RuneCountInString(task) == len(task) {
		len_task += len(task)
	}
	format := fmt.Sprintf("%%-2d %%-%ds %%-%ds %%s\n", len_status, len_task)
	return fmt.Sprintf(format, id, status, task, limit.Format("2006/01/02 15:04:05"))
}

func list() {
	todos, _ := loadTodos()
	if len(todos) == 0 {
		fmt.Println("📭 Todoはありません")
		return
	}
	fmt.Printf("%-2s %-7s %-21s %-19s\n", "ID", "ステータス", "タスク", "期限")
	fmt.Printf("%-2s %-12s %-24s %-19s\n", "--", "------------", "------------------------", "-------------------")
	for _, t := range todos {
		status := "未完了"
		if t.IsDone {
			status = "完了"
		}
		if t.IsDeleted {
			status = "削除済"
		}
		fmt.Print(format(t.ID, status, t.Task, t.Limit))
	}
}

func done(id int) {
	todos, _ := loadTodos()
	for i, t := range todos {
		if t.ID == id && !t.IsDone {
			todos[i].Done = true
			todos[i].IsDone = true
			saveTodos(todos)
			fmt.Println("✅ 完了:", t.Task)
			return
		} else if t.ID == id {
			fmt.Println("⚠️ すでに完了しています")
			return
		}
	}
	fmt.Println("⚠️ IDが見つかりません")
}

func deleteTodo(id int) {
	todos, _ := loadTodos()
	for i, t := range todos {
		if t.ID == id && t.IsDone {
			fmt.Println("⚠️ 完了しているタスクは削除できません")
			return
		} else if t.ID == id && t.IsDeleted {
			fmt.Println("⚠️ すでに削除されています")
			return
		} else if t.ID == id {
			todos[i].IsDeleted = true
			todos[i].DeletedAt = time.Now()
			saveTodos(todos)
			fmt.Println("🗑️ 削除:", t.Task)
			return
		}
	}
	fmt.Println("⚠️ IDが見つかりません")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("📘 Todo CLI アプリ - 'help' または 'h' でコマンド一覧")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		switch command {
		case "help", "h":
			fmt.Println("使用可能なコマンド:")
			fmt.Println("  add(a) <タスク名>      - タスクを追加")
			fmt.Println("  list(l)               - タスク一覧表示")
			fmt.Println("  done(d) <ID>          - 指定IDのタスクを完了")
			fmt.Println("  delete(del) <ID>      - 指定IDのタスクを削除")
			fmt.Println("  exit(e)               - 終了")

		case "add", "a":
			if len(args) == 0 {
				fmt.Println("⚠️ タスクを入力してください")
				continue
			}
			task := strings.Join(args, " ")
			add(task, reader)

		case "list", "l":
			list()

		case "done", "d":
			if len(args) == 0 {
				fmt.Println("⚠️ IDを入力してください")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("⚠️ 数値のIDを入力してください")
				continue
			}
			done(id)

		case "delete", "del":
			if len(args) == 0 {
				fmt.Println("⚠️ IDを入力してください")
				continue
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("⚠️ 数値のIDを入力してください")
				continue
			}
			deleteTodo(id)

		case "exit", "e":
			fmt.Println("👋 終了します")
			return

		default:
			fmt.Println("❓ 未知のコマンド:", command)
		}
	}
}
