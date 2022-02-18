import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import sys
import os

def load_data(csv_path):
    data = pd.read_csv(csv_path)
    data['throughput'] = 1000 * data['total requests'] / data['latency in ms']  # requests per second
    requests_df = pd.concat(
        [
            data.groupby('total requests')['latency in ms'].mean(),
            data.groupby('total requests')['throughput'].mean()
        ], 
        axis=1
    )

    clients_df = pd.concat(
        [
            data.groupby('clients')['latency in ms'].mean(),
            data.groupby('clients')['throughput'].mean()
        ], 
        axis=1
    )

    return clients_df, requests_df

def save_figures(clients_df, requests_df, prefix):
    s = 5
    ysize = 0.8 * s
    xsize = 2 * s

    # Latency Plots
    plt.figure(figsize=(xsize, ysize))
    plt.suptitle('Latency')

    plt.subplot(121)
    plt.plot(clients_df.index, clients_df['latency in ms'], 'ro')
    plt.plot(np.unique(clients_df.index), np.poly1d(np.polyfit(clients_df.index, clients_df['latency in ms'], 1))(np.unique(clients_df.index)), 'b--')
    plt.ylabel('Latency (in ms)')
    plt.xlabel('No. of clients')

    plt.subplot(122)
    plt.plot(requests_df.index, requests_df['latency in ms'], 'bo')
    plt.plot(np.unique(requests_df.index), np.poly1d(np.polyfit(requests_df.index, requests_df['latency in ms'], 1))(np.unique(requests_df.index)), 'r--')
    plt.ylabel('Latency (in ms)')
    plt.xlabel('Total No. of requests')
    plt.savefig(f'{prefix}_latency.png')

    # Throughput Plots
    plt.figure(figsize=(xsize, ysize))
    plt.suptitle('Throughput')

    plt.subplot(121)
    plt.plot(clients_df.index, clients_df['throughput'], 'ro')
    plt.plot(np.unique(clients_df.index), np.poly1d(np.polyfit(clients_df.index, clients_df['throughput'], 1))(np.unique(clients_df.index)), 'b--')
    plt.ylabel('Throughput (requests per second)')
    plt.xlabel('No. of clients')

    plt.subplot(122)
    plt.plot(requests_df.index, requests_df['throughput'], 'bo')
    plt.plot(np.unique(requests_df.index), np.poly1d(np.polyfit(requests_df.index, requests_df['throughput'], 1))(np.unique(requests_df.index)), 'r--')
    plt.ylabel('Throughput (requests per second)')
    plt.xlabel('Total No. of requests')
    plt.savefig(f'{prefix}_throughput.png')


if __name__ == '__main__':
    if len(sys.argv) >= 2:
        csv_path = sys.argv[1]
    else:
        sys.exit(1)  # Not enough arguments provided

    # Parse the data
    clients_df, requests_df = load_data(csv_path)

    # Save the figures
    prefix = os.path.splitext(os.path.basename((csv_path)))[0]
    save_figures(clients_df, requests_df, prefix)
