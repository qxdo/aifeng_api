package common

type Person struct {
	Proxy string
	Guid  string
}

type Config struct {
	MaxThread int `json:"maxThread"`
	TaskNum   int `json:"taskNum"`
	Database  Database
}

type Database struct {
	Driver     string `json:"driver"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Datasource string `json:"datasource"`
}

// Task 任务实体
type Task struct {
	Id          int     `json:"id" db:"id"`
	Uid         int     `json:"uid" db:"uid"`
	ReadUrl     string  `json:"readUrl,omitempty"`
	Url         string  `json:"url,omitempty" db:"url"`
	Guid        string  `json:"guid,omitempty" db:"guid"`
	Title       string  `json:"title,omitempty" db:"title"`
	Type        string  `json:"type,omitempty" db:"type"`
	Price       float32 `json:"price" db:"price"`
	DemandCount int     `json:"demand_count,omitempty" db:"demand_count"`
	BeforeCount int     `json:"before_count" db:"before_count"`
	AllCount    int     `json:"all_count" db:"all_count"`
	SucCount    int     `json:"suc_count" db:"suc_count"`
	Status      string  `json:"status,omitempty" db:"status"`
	AddTime     string  `json:"add_time,omitempty" db:"add_time"`
	EndTime     string  `json:"end_time,omitempty" db:"end_time"`
	IsFirst     int     `json:"is_first" db:"is_first"`
	Start       string  `json:"start,omitempty" db:"start"`
	Secret      string  `json:"secret,omitempty" db:"secret"`
	Priority    int     `json:"priority" db:"priority"`
	IsDelete    int     `json:"is_delete,omitempty" db:"is_delete"`
}

type Params struct {
	Url string `json:"url"`
	Num int    `json:"num"`
}

type ProxyEntry struct {
	Id    int    `json:"id" db:"id"`
	Proxy string `json:"proxy" db:"proxy"`
	Guid  string `json:"guid" db:"guid"`
	Count int    `json:"count" db:"count"`
	Time  string `json:"time" db:"time"`
}

// CountInfo 前端显示数量信息实体
type CountInfo struct {
	AllAccount     int     `json:"all_account" db:"all_account"`
	EnableAccount  int     `json:"enable_account" db:"enable_account"`
	SleepAccount   int     `json:"sleep_account" db:"sleep_account"`
	TaskNum        int     `json:"task_num" db:"task_num"`
	RunningNum     int     `json:"running_num" db:"running_num"`
	CompletedNum   int     `json:"completed_num" db:"completed_num"`
	TaskCount      int     `json:"task_count" db:"task_count"`
	CompletedCount int     `json:"completed_count" db:"completed_count"`
	TotalPrice     float64 `json:"total_price" db:"total_price"`
	DayCount       int     `json:"day_count" db:"day_count"`
	NightCount     int     `json:"night_count" db:"night_count"`
	YestNightCount int     `json:"yest_night_count" db:"yest_night_count"`
	YestDayCount   int     `json:"yest_day_count" db:"yest_day_count"`
}

// User 用户实体
type User struct {
	Id       int    `json:"_" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Disabled int    `json:"disabled" db:"disabled"`
}

// Price 价单
type Price struct {
	Id         int     `json:"_" db:"id"`
	DayPrice   float64 `json:"day_price" db:"day_price"`
	NightPrice float64 `json:"night_price" db:"night_price"`
	Uid        int     `json:"uid" db:"uid"`
	SqlTime    string  `json:"sqlTime" db:"time"`
}
