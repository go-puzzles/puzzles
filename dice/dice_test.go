// File:		dice_test.go
// Created by:	Hoven
// Created on:	2025-03-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dice

import (
	"math"
	"testing"
)

func TestDice_Next(t *testing.T) {
	tests := []struct {
		name    string
		weights []int
	}{
		{
			name:    "正常权重",
			weights: []int{1, 2, 3},
		},
		{
			name:    "全零权重",
			weights: []int{0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDice(tt.weights)
			seen := make(map[int]bool)

			// 对于非零权重的测试，确保所有索引都能被抽到
			if tt.name == "正常权重" {
				for i := 0; i < len(tt.weights); i++ {
					result := d.Next()
					if result == -1 {
						t.Errorf("未期望返回 -1")
					}
					seen[result] = true
				}

				// 确保所有可能的索引都被抽到过
				for i := range tt.weights {
					if !seen[i] {
						t.Errorf("索引 %d 未被抽到", i)
					}
				}

				// 权重用尽后应返回 -1
				if result := d.Next(); result != -1 {
					t.Errorf("权重用尽后期望返回 -1，实际返回 %d", result)
				}
			}

			// 对于全零权重的测试
			if tt.name == "全零权重" {
				if result := d.Next(); result != -1 {
					t.Errorf("全零权重时期望返回 -1，实际返回 %d", result)
				}
			}
		})
	}
}

func TestDice_Reset(t *testing.T) {
	weights := []int{1, 2, 3}
	d := NewDice(weights)

	// 先抽取一些值
	d.Next()
	d.Next()

	// 重置
	d.Reset()

	if d.total != 6 {
		t.Errorf("重置后期望 total 为 6，实际得到 %d", d.total)
	}

	// 验证重置后可以继续抽取
	seen := make(map[int]bool)
	for i := 0; i < len(weights); i++ {
		result := d.Next()
		if result == -1 {
			t.Errorf("重置后未期望返回 -1")
		}
		seen[result] = true
	}

	// 确保所有可能的索引都被抽到过
	for i := range weights {
		if !seen[i] {
			t.Errorf("重置后索引 %d 未被抽到", i)
		}
	}
}

func TestDice_NextRate(t *testing.T) {
	times := 10000
	weights := []int{1, 2, 8}
	cnt := make(map[int]int)

	dice := NewDice(weights)
	for i := 0; i < times; i++ {
		n := dice.Next()
		dice.Reset()

		cnt[n]++
	}

	t.Log(cnt)
	for i := 0; i < len(weights); i++ {
		cnt[i] /= weights[i]
	}
	t.Log(cnt)

	if cnt[-1] > 0 {
		t.Error("Got -1 from Next(), expect no -1")
	}
	for i := 1; i < len(weights); i++ {
		if math.Abs(float64(cnt[i]-cnt[i-1])) >= float64(times)*0.01 {
			t.Error("The result probability is not equals to weights within error rate 1%")
		}
	}
}

func TestDice_DrawingOrder(t *testing.T) {
	// 设置权重：最后一个数权重最大，理论上最容易先被抽到
	weights := []int{1, 3, 8, 20}
	d := NewDice(weights)

	// 记录每个位置作为第一个被抽到的数字的次数
	firstDrawCounts := make([]int, len(weights))
	totalRounds := 10000 // 测试轮数

	for round := 0; round < totalRounds; round++ {
		// 记录第一次抽取的结果
		firstDraw := d.Next()
		firstDrawCounts[firstDraw]++

		// 完成这一轮的剩余抽取
		for {
			n := d.Next()
			if n == -1 {
				d.Reset()
				break
			}
		}
	}

	// 验证权重大的数字更容易先被抽到
	t.Logf("首次抽取统计（总轮数：%d）：", totalRounds)
	for i, count := range firstDrawCounts {
		percentage := float64(count) / float64(totalRounds) * 100
		t.Logf("索引 %d (权重 %d): %d 次 (%.2f%%)", i, weights[i], count, percentage)
	}

	// 验证权重最大的位置（索引3）应该最常作为第一个被抽到的数
	if firstDrawCounts[3] <= firstDrawCounts[2] ||
		firstDrawCounts[3] <= firstDrawCounts[1] ||
		firstDrawCounts[3] <= firstDrawCounts[0] {
		t.Errorf("权重分布异常：最大权重位置 (索引3) 首次抽取次数不是最多的")
	}

	// 验证权重顺序关系
	if firstDrawCounts[2] <= firstDrawCounts[1] ||
		firstDrawCounts[1] <= firstDrawCounts[0] {
		t.Errorf("权重分布异常：首次抽取次数未按权重大小排序")
	}
}
