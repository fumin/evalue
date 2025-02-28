# This script visualizes the output of TestOptionalContinuation.

import numpy as np
import matplotlib
import matplotlib.pyplot as plt
import pandas as pd


def plotP(df):
  numSamples = len(df.s.unique())
  np.random.seed(0)
  plottedSamples = np.random.choice(range(numSamples), 200, replace=False)
  xmin, xmax = df.t.min(), df.t.max()
  
  fig, ax = plt.subplots(1, 1)
  for sampleID in plottedSamples:
    sdf = df[df.s == sampleID]
    row = sdf.iloc[0]

    if row.stopTP != -1:
      color = (0, 1, 0)
      alpha = 1
      ax.plot(sdf.t, sdf.p, color=color, alpha=alpha)
    elif row.stopTPOC != -1:
      sdfReal = sdf[sdf.t <= row.stopTPOC]
      sdfHypothetic = sdf[sdf.t > row.stopTPOC]

      color = (1, 0, 0)
      alpha = 1
      ax.plot(sdfReal.t, sdfReal.p, color=color, alpha=alpha)
      ax.plot(sdfHypothetic.t, sdfHypothetic.p, linestyle="--", color=color, alpha=alpha)
    else:
      color = (0, 0, 0)
      alpha = 0.1
      ax.plot(sdf.t, sdf.p, color=color, alpha=alpha)

  ax.hlines(0.05, xmin, xmax, linewidth=3)

  handles = [
      [matplotlib.patches.Patch(facecolor=(0, 1, 0))],
      [matplotlib.patches.Patch(facecolor=(0, 1, 0)), matplotlib.patches.Patch(facecolor=(1, 0, 0))],
      ]
  labels = [
      "standard, type I error {:.3f}".format(len(df[df.stopTP != -1].s.unique())/numSamples),
      "optional continuation, type I error {:.3f}".format(len(df[df.stopTPOC != -1].s.unique())/numSamples),
      ]
  ax.legend(handles=handles, labels=labels, handler_map={list: matplotlib.legend_handler.HandlerTuple(None)})

  ax.set_yscale('log')
  ax.set_xlabel("n")
  ax.set_ylabel("p-value")
  ax.set_title("T-test with/without optional continuation")


def plotE(df):
  stoppedStd = list(df[df.stopTE != -1].s.unique())
  stoppedOC = list(df[(df.stopTEOC != -1) & (df.stopTE == -1)].s.unique())

  numSamples = len(df.s.unique())
  np.random.seed(0)
  numPlotted = 200
  plottedSamples = list(np.random.choice(range(numSamples), numPlotted, replace=False))
  plottedOC = list(np.random.choice(stoppedOC, int(np.round(len(stoppedOC)*numPlotted/numSamples)), replace=False))
  # plottedSamples += stoppedStd + plottedOC
  xmin, xmax = df.t.min(), df.t.max()
  
  fig, ax = plt.subplots(1, 1)
  for sampleID in plottedSamples:
    sdf = df[df.s == sampleID]
    row = sdf.iloc[0]

    if row.stopTE != -1:
      color = (0, 1, 0)
      alpha = 1
      ax.plot(sdf.t, sdf.e, color=color, alpha=alpha)
    elif row.stopTEOC != -1:
      sdfReal = sdf[sdf.t <= row.stopTEOC]
      sdfHypothetic = sdf[sdf.t > row.stopTEOC]

      color = (1, 0, 0)
      alpha = 1
      ax.plot(sdfReal.t, sdfReal.e, color=color, alpha=alpha)
      ax.plot(sdfHypothetic.t, sdfHypothetic.e, linestyle="--", color=color, alpha=alpha)
    else:
      color = (0, 0, 0)
      alpha = 0.1
      ax.plot(sdf.t, sdf.e, color=color, alpha=alpha)

  ax.hlines(1./0.05, xmin, xmax, linewidth=3)

  handles = [
      [matplotlib.patches.Patch(facecolor=(0, 1, 0))],
      [matplotlib.patches.Patch(facecolor=(0, 1, 0)), matplotlib.patches.Patch(facecolor=(1, 0, 0))],
      ]
  labels = [
      "standard, type I error {:.3f}".format(len(df[df.stopTE != -1].s.unique())/numSamples),
      "optional continuation, type I error {:.3f}".format(len(df[df.stopTEOC != -1].s.unique())/numSamples),
      ]
  ax.legend(handles=handles, labels=labels, handler_map={list: matplotlib.legend_handler.HandlerTuple(None)})

  ax.set_yscale('log')
  ax.set_xlabel("n")
  ax.set_ylabel("e-value")
  ax.set_title("E-value test with/without optional continuation")


df = pd.read_csv("oc.csv")
plotP(df)
plotE(df)
