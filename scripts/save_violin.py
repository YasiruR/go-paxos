import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import sys
import os

def load_data(csv_path):
    data = pd.read_csv(csv_path)
    data['throughput'] = 1000 * data['total requests'] / data['latency in ms']  # requests per second

    return data

def save_figures(data, prefix):
    sns.violinplot(x='total requests', y='throughput', data=data, palette='Set3', scale='width'); plt.title('Throughput (requests per second)'); plt.savefig(f'{prefix}_throughput_violin.png')
    # sns.violinplot(x='total requests', y='latency in ms', data=data, palette='Set3', scale='width'); plt.title('Latency (ms)'); plt.savefig(f'{prefix}_violin.png')

if __name__ == '__main__':
    if len(sys.argv) >= 2:
        csv_path = sys.argv[1]
    else:
        sys.exit(1)  # Not enough arguments provided

    # Parse the data
    data = load_data(csv_path)

    # Save the figures
    prefix = os.path.splitext(os.path.basename((csv_path)))[0]
    save_figures(data, prefix)
