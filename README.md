# HyperLedger-Fabric-Smart-Contract-v1


This repository contains a Hyperledger Fabric chaincode written in Go. The chaincode provides functionality for managing trainers, trainees, platforms, and virtual labs.

## Prerequisites

- Go 1.14 or later
- Hyperledger Fabric 2.0 or later

## Installation

1. Clone the repository:


2. Build the chaincode:


## Usage

The chaincode provides the following functions:

- `createTrainee`: Creates a new trainee with the specified details.
- `createPlatform`: Creates a new platform with the specified details.
- `addTraineeToPlatform`: Adds a trainee to a platform.
- `delete`: Deletes an entity from the ledger.
- `ScoreTheVlab`: Scores a virtual lab for a trainee.
- `createTrainer`: Creates a new trainer with the specified details.
- `getIdentity`: Retrieves the identity (trainee/trainer) based on the provided ID.

To deploy the chaincode, follow the instructions provided by the Hyperledger Fabric documentation.

## Examples

### Creating a Trainee

