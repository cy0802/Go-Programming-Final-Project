package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Event struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type dayDiary struct {
	Date   string  `json:"date"`
	Events []Event `json:"events"`
}

var noteBook = []dayDiary{}
var currentId = 1 // 用來生成唯一 ID

// 如果點擊到已經POST過的日期，可以用GET請求直接得到之前存好的值
func getNotebook(c *gin.Context) {
	if len(noteBook) > 0 {
		c.JSON(200, noteBook) // 回傳當天所有事件
		return
	}
	c.JSON(404, gin.H{"message": "No events found in the notebook"})
}

// 已經POST過的特定Event可以對其進行修改
func updateEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("eventId")) // 獲取事件的 Id
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid event ID"})
		return
	}
	var updatedEvent Event
	err = c.ShouldBindJSON(&updatedEvent)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}
	updatedEvent.Id = eventId
	//找到指定Id的Event並修改它
	for i, _ := range noteBook {
		for j, event := range noteBook[i].Events {
			if event.Id == eventId {
				noteBook[i].Events[j].Title = updatedEvent.Title
				noteBook[i].Events[j].Content = updatedEvent.Content
				c.JSON(200, gin.H{"message": "Event updated", "event": noteBook[i].Events[j]}) //回傳單一事件
				return
			}
		}
	}
	c.JSON(404, gin.H{"error": "Event not found"})
}

// 刪除一個指定Id的Event
func deleteEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("eventId")) // 獲取事件的 Id
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid event ID"})
		return
	}
	for i, day := range noteBook {
		for j, event := range noteBook[i].Events {
			if event.Id == eventId {
				noteBook[i].Events = append(day.Events[:j], day.Events[j+1:]...)
				c.JSON(200, gin.H{"message": "Event deleted"})
				return
			}
		}
	}
	c.JSON(404, gin.H{"error": "Event not found"})
}

// POST一個Event的資訊(Title+Time+Content)儲存在Server中，並給予每一個Event一個獨立的Id
func addEvent(c *gin.Context) {
	date := c.Param("date")
	var newEvent Event
	err := c.ShouldBindJSON(&newEvent)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}
	newEvent.Id = currentId
	currentId++
	//如果該日期已存在，直接追加事件
	for i, day := range noteBook {
		if day.Date == date {
			noteBook[i].Events = append(noteBook[i].Events, newEvent)
			c.JSON(200, gin.H{"message": "Event added", "date": date, "event": newEvent}) //回傳單一事件
			return
		}
	}
	//如果該日期不存在，新增日期並追加事件
	newDay := dayDiary{Date: date, Events: []Event{newEvent}}
	noteBook = append(noteBook, newDay)                                                       //儲存Date+Content
	c.JSON(200, gin.H{"message": "New day and event added", "date": date, "event": newEvent}) //回傳單一事件
}

// 將選定範圍內的日期的所有事件都丟給串聯AI的API函數，並回傳結果給前端
func summarizeWeek(c *gin.Context) {
	startDate := c.Param("startDate")
	endDate := c.Param("endDate")
	start, err1 := time.Parse("2006-01-02", startDate)
	end, err2 := time.Parse("2006-01-02", endDate)
	if err1 != nil || err2 != nil {
		c.JSON(409, gin.H{"message": "Invalid date format"})
		return
	}
	//收集這一整週的日記
	var weekDiary = []dayDiary{}
	for _, day := range noteBook {
		dayDate, _ := time.Parse("2006-01-02", day.Date)
		if (dayDate.Equal(start) || dayDate.After(start)) && dayDate.Before(end.AddDate(0, 0, 1)) {
			weekDiary = append(weekDiary, day)
		}
	}
	var summaryBuilder strings.Builder
	summaryBuilder.WriteString("本週總結：")
	for _, day := range weekDiary {
		dayDate, _ := time.Parse("2006-01-02", day.Date)
		weekday := dayDate.Weekday()
		summaryBuilder.WriteString(weekday.String() + ": ")
		for j, event := range day.Events {
			eventIdInADay := j + 1
			summaryBuilder.WriteString("Event" + strconv.Itoa(eventIdInADay) + ":")
			summaryBuilder.WriteString("{title:" + event.Title + ";content:" + event.Content + "}")
			eventIdInADay++
		}
	}
	summary := summaryBuilder.String()
	//傳本週總結給AI
	aiResponse := callAI(summary)
	c.JSON(200, gin.H{"summary": aiResponse}) //回傳AI生成的內容
}

func main() {
	router := gin.Default()
	// Enable CORS
	router.Use(cors.Default())
	router.RedirectFixedPath = true
	router.GET("/diary", getNotebook)                         //需要一個帶有日期的URL (日期格式:"YYYY-MM-DD")
	router.PUT("/diary/:eventId", updateEvent)                //需要一個帶有指定事件Id的URL; 傳來的JSON需要包含新的Title+Time+Cntent
	router.DELETE("/diary/:eventId", deleteEvent)             //需要一個帶有指定事件Id的URL
	router.POST("/diary/:date", addEvent)                     //需要一個帶有日期的URL;傳來的JSON需要包含Title+Time+Cntent
	router.GET("/summary/:startDate/:endDate", summarizeWeek) //前端需要先確認好該週的範圍(起始日和終止日)
	err := router.Run(":8083")
	if err != nil {
		return
	}
}
