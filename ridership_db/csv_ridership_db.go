package ridershipDB

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type CsvRidershipDB struct {
	idIdxMap      map[string]int
	csvFile       *os.File
	csvReader     *csv.Reader
	num_intervals int
}

func (c *CsvRidershipDB) Open(filePath string) error {
	c.num_intervals = 9

	// Create a map that maps MBTA's time period ids to indexes in the slice
	c.idIdxMap = make(map[string]int)
	for i := 1; i <= c.num_intervals; i++ {
		timePeriodID := fmt.Sprintf("time_period_%02d", i)
		c.idIdxMap[timePeriodID] = i - 1
	}

	// create csv reader
	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	c.csvFile = csvFile
	c.csvReader = csv.NewReader(c.csvFile)

	return nil
}

// TODO: some code goes here
// Implement the remaining RidershipDB methods
func (c *CsvRidershipDB) GetRidership(lineId string) ([]int64, error) {
	// 初始化一个长度为 9 的切片，用于存储每个时间段的乘客总数
	ridership := make([]int64, c.num_intervals)

	// 重新定位到文件开始，以便重新读取
	_, err := c.csvFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	c.csvReader = csv.NewReader(c.csvFile)

	// 忽略 CSV 文件的第一行（通常是标题行）
	_, err = c.csvReader.Read()
	if err != nil {
		return nil, err
	}

	// 遍历 CSV 文件的每一行
	for {
		record, err := c.csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// 假设 record 中包含了线路 ID、时间段 ID 和乘客数量
		// 且我们知道它们的列位置
		currentLineId := record[0]                                 // 这里需要知道线路 ID 在哪一列
		timePeriodId := record[2]                                  // 这里需要知道时间段 ID 在哪一列
		ridershipCount, err := strconv.ParseInt(record[4], 10, 64) // 这里需要知道乘客数量在哪一列
		if err != nil {
			return nil, err
		}

		// 检查线路是否匹配
		if currentLineId == lineId {
			// 获取时间段的索引
			timePeriodIdx, exists := c.idIdxMap[timePeriodId]
			if exists {
				// 累加乘客数量
				ridership[timePeriodIdx] += ridershipCount
			}
		}
	}

	return ridership, nil
}

func (c *CsvRidershipDB) Close() error {
	return nil
}
