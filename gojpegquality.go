package gojpegquality

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"time"
)

func searchQuantizationTable(buf []byte) [3]string {
	var (
		B1     byte
		B2     byte
		length int
		arr    [3]string
		table  []byte
		typer  int
		i      int
	)

	for i = 2; i < len(buf); {
		if buf[i] == 0xff {
			B1 = buf[i]
			i += 1
			for buf[i] == 0xff {
				i += 1
			}
			if buf[i] == 0 {
				continue
			}
			B2 = buf[i]
			i += 1
			if i >= len(buf) {
				break
			}
			typer = int(B1)*256 + int(B2)
			length = int(buf[i])*256 + int(buf[i+1]) - 2
			i += 2
			if typer != 0xffdb {
				i += length
				continue
			}
			if length%65 == 0 {
				table = buf[i : i+length]
				for j := 0; j < length/65; j += 1 {
					arr[j] = hex.EncodeToString(table[j*65 : (j+1)*65])
				}
				i += length
			}
		} else {
			i += 1
		}
	}
	return arr
}

func averageTable(table string, index int) float64 {
	var arr []string
	var total float64
	var dd int16

	for i := 0; i < len(table); i += 2 {
		arr = append(arr, string(table[i])+string(table[i+1]))
	}

	for j := 0; j < len(arr); j += 1 {
		d, _ := strconv.ParseInt(arr[j], 16, 64)
		dd = int16(d)
		if j != 0 {
			total += float64(dd)
		}
	}

	result := 100.00 - total/63
	return result
}

func GetQ(buf []byte) float64 {
	start := time.Now()
	var avgs [3]float64

	data := buf

	if data[0] != 0xFF || data[1] != 0xD8 {
		fmt.Println("ERROR: Not a supported JPEG format")
		return 1
	}

	tables := searchQuantizationTable(data)

	for key, value := range tables {
		if value != "" {
			avgs[key] = averageTable(value, key)
		}
	}

	if len(avgs) == 2 || avgs[2] == 0 {
		avgs[2] = avgs[1]
		diff := math.Abs(avgs[0]-avgs[1]) * 0.49
		diff += math.Abs(avgs[0]-avgs[2]) * 0.49
		quality := (avgs[0]+avgs[1]+avgs[2])/3 + diff
		fmt.Println("avgs == 2 ->", quality)
		elapsed := time.Since(start)
		fmt.Println("time:", elapsed)
		return quality
	}

	if len(avgs) == 1 || (avgs[1] == 0 && avgs[2] == 0) {
		fmt.Println(avgs)
		return avgs[0]
	}

	if len(avgs) == 3 {
		diff := math.Abs(avgs[0]-avgs[1]) * 0.49
		diff += math.Abs(avgs[0]-avgs[2]) * 0.49
		quality := (avgs[0]+avgs[1]+avgs[2])/3 + diff
		fmt.Println("avgs == 3 ->", quality)
		return quality
	}

	return 1
}
