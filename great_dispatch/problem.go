package main

/*

Problem

Given N boxes Bj of weight and volume Wj, Vj
and M trucks Ti of capacity Ci (max = 100)

Need to put boxes in ALL trucks

A_ij = 1 if Bj is in Ti else 0

Find a solution such as:
 - Each truck has a volume max of 100: for each truck Ti: Sum Aij * Vj <= 100
 - Each box is one truck only: for each box Bj: Sum Aij = 1
 - The difference between the heavier and lighter truck is minimal:
   Max(Sum Wj) - Min(Sum Wj) is minimal

Given i : Sum Aij * Vj + Si = 100
Given j : Sum Aij = 1
Argmin_i Sum_j(AijWj) - Argmax_i(AijWj)


*/
