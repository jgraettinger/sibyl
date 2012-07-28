package cclparse

import (
    "sort"
)

type LabelWeights map[Label]float64

type LabelWeightsEntry struct {
    Label Label
    Weight float64
}

type LabelWeightsList []LabelWeightsEntry

func (list LabelWeightsList) Len() int {
    return len(list)
}
func (list LabelWeightsList) Less(i, j int) bool {
    if list[i].Weight == list[j].Weight {
        return list[i].Label.Token < list[j].Label.Token
    }
    return list[i].Weight > list[j].Weight
}
func (list LabelWeightsList) Swap(i, j int) {
    list[i], list[j] = list[j], list[i]
}

func (weights LabelWeights) FilterToTopN(count int) LabelWeights {

    var list LabelWeightsList

    for label, weight := range(weights) {
        list = append(list, LabelWeightsEntry{label, weight})
    }
    sort.Sort(list)

    filteredWeights := make(LabelWeights)
    for index, entry := range(list) {
        if index == count {
            break
        }
        filteredWeights[entry.Label] = entry.Weight
    }
    return filteredWeights
}

