#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Feb 15 20:26:24 2022

@author: yasi
"""

import csv
from matplotlib import pyplot as plt

def read():
    clients = []
    reqs = []
    total_reqs = []
    success_reqs = []
    latency = []
    with open('results.csv') as file:
        reader = csv.reader(file)
        line_index = 0
        for row in reader:
            if line_index > 1:
                clients.append(int(row[0]))
                reqs.append(int(row[1]))
                total_reqs.append(int(row[2]))
                success_reqs.append(int(row[3]))
                latency.append(int(row[4]))
            line_index += 1
    return clients, reqs, total_reqs, success_reqs, latency

def plot(clients, reqs, total, success_reqs, latency):
    rates= []
    for i in range(len(total)):
        rates.append(float( latency[i]/(total[i]*1000) ))

    print(clients, rates)

    plt.plot(clients, rates, label='actual')
    plt.legend(loc='upper left')
    plt.xlabel('number of concurrent clients')
    plt.ylabel('requests per second')
    plt.savefig('rate_vs_clients.png')
    plt.show()

clients, reqs, total_reqs, success_reqs, latency = read()
plot(clients, reqs, total_reqs, success_reqs, latency)