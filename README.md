# GO DEEP

## Motivation

This is my personal education project. Things I want to train and learn with this project: AI deep learning basics, Go IoC architecture approaches, Go programs performance optimization and profiling. The goal is to assemble a flexible and efficient deep learning framework on pure Go.

## Installation

`go get github.com/I159/go_deep`

And if you need examples to play with there is a sample project with examples for the library.

`go get github.com/I159/go_deep_examples`

## Contribution

If you have the same goals of learning or/and you have more solid math or architectural background than me feel free to fork and make pull requests.

## Current state of the project

I spent a lot of time to build architecturally correct project but I lost algorithm itself. Perceptron (I started from this simplest NN in my research) could be implemented in 50 lines of code easily but with all my layers of abstraction and overabundant asynchrony I drawned in massive debug. Then I've been busy with billed projects and completly lost the point where the things went wrong. I spent several hours to recreate backpropagation in strightforward manner and then grow layers and abltraction on workable algorithm mostly from scratch.