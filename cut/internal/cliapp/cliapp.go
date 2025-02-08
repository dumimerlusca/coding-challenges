package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type CliApp struct {
	cfg    Config
	r      io.Reader
	output *bytes.Buffer
}

func NewApp(cfg Config, r io.Reader) *CliApp {
	return &CliApp{cfg: cfg, r: r, output: bytes.NewBuffer([]byte{})}
}

func (app *CliApp) Run() ([]byte, error) {

	switch app.cfg.Select.OptionType {
	case SelectOptionFields:
		return app.selectFields()
	default:
		return nil, fmt.Errorf("unsupported select option: %v", app.cfg.Select.OptionType)
	}

}

func (app *CliApp) selectFields() ([]byte, error) {
	selectedFields, err := parseSelectValue(app.cfg.Select.Value)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(app.r)

	for {
		line := scanner.Text()
		fieldNumber := 1
		start := 0
		end := 0
		isFirstField := true

		for idx, char := range line {
			if string(char) == app.cfg.Delimeter {
				start = end
				end = idx

				if !slices.Contains(selectedFields, fieldNumber) {
					fieldNumber++
					continue
				}
				fieldNumber++

				var s string

				if isFirstField {
					s = strings.TrimLeft(line[start:end], app.cfg.Delimeter)
				} else {
					s = line[start:end]
				}

				isFirstField = false

				_, err := app.output.WriteString(s)
				if err != nil {
					return nil, err
				}

			}
		}

		if scanner.Scan() {
			_, err = app.output.WriteRune('\n')
			if err != nil {
				return nil, err
			}

		} else {
			break
		}

	}

	return app.output.Bytes(), nil
}

func parseSelectValue(value string) ([]int, error) {

	res := []int{}

	groups := strings.Split(value, ",")

	for _, group := range groups {
		if strings.Contains(group, "-") {
			groupLimits := strings.Split(group, "-")
			if len(groupLimits) > 2 {
				return nil, fmt.Errorf("invalid interval: %v", group)
			}
			left, err := strconv.Atoi(groupLimits[0])
			if err != nil {
				return nil, fmt.Errorf("invalid interval: %v", group)
			}

			right, err := strconv.Atoi(groupLimits[1])
			if err != nil {
				return nil, fmt.Errorf("invalid interval: %v", group)
			}

			if left > right {
				return nil, fmt.Errorf("invalid interval: %v", group)
			}

			for i := left; i <= right; i++ {
				res = append(res, i)
			}
		} else {
			val, err := strconv.Atoi(group)
			if err != nil {
				return nil, fmt.Errorf("invalid value: %v", group)
			}

			res = append(res, val)

		}
	}

	return res, nil
}
