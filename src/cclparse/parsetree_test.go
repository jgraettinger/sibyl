package cclparse

import (
    "testing"
)


func TestExtractConstituentParse(t * testing.T) {
    chart := NewChart([]string{"I", "know", "the", "boy", "sleeps"})

    chart.AddLink(1, 0, 1) // know =1> I
    chart.AddLink(1, 2, 0) // know =0> the
    chart.AddLink(2, 3, 0) // the =0> boy
    chart.AddLink(3, 2, 0) // boy =0> the
    chart.AddLink(1, 4, 0) // know =0> sleeps
    chart.AddLink(4, 3, 1) // sleeps =1> boy

    t.Log(chart)

    heights := ExtractSyntacticParse(chart)

    t.Log(heights)

}

/*
    chart := NewChart([]string{"I", "heard", "a", "man", "rode", "the", "subway", "on", "friday"})

    chart.AddLink(1, 0, 1) // heard =1> I
    chart.AddLink(1, 2, 0) // heard =0> a
    chart.AddLink(2, 3, 0) // a =0> man
    chart.AddLink(3, 2, 0) // man =0> a
    chart.AddLink(1, 4, 0) // heard =0> rode
    chart.AddLink(4, 3, 1) // rode =1> man
    chart.AddLink(4, 5, 0) // rode =0> the
    chart.AddLink(5, 6, 0) // the =0> subway
    chart.AddLink(6, 5, 0) // subway =0> the
    chart.AddLink(6, 7, 0) // subway =0> on
    chart.AddLink(7, 8, 0) // on =0> friday
    chart.AddLink(8, 7, 0) // friday =0> on
*/
