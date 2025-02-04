# VM to container migration system

This is the GitHub repo for the project "Virtual machine to container migration system", created as a research project during semester 1 of my Master in Applied IT at Fontys University. The repository showcase a PoC for a system that can perform automated migrations of stateless applications such web servers from Ubuntu virtual machines to Docker containers.

## Repository structure

### Structure overview

```sh
vm-to-container-migrator/
│
├── api/             # The "api" folder contains the code base for the back-end API part of the project, 
│                    # which serves the function of a core unit handling the business logic to perform 
│                    # VM-to-container automated migrations.
│
├── cli/             # The "cli" folder includes the code for the user interface of the system, 
│                    # implemented in the form of a CLI with which the user can interact to ask certain 
│                    # applications to be migrated from a virtual machine to a Docker container.
│
├── experiments/     # The "experiments" folder contains various experiments developed during the project 
                     # to test different parts of the migration process and also learn before assembling 
                     # the final version of the PoC.
```

### API

```sh
api/
│
├── internal/         # It is a magic folder in Golang that prevents other projects from importing code 
│   │                 # from this directory. The following folders are placed here:
│   │
│   ├── model/        # Contains the models, a blueprint for the data structures that are part of the 
│   │                 # project's database.
│   │
│   ├── route/        # Contains API endpoints, each associated with a function executed when a request 
│   │   │             # is made. Also includes the service/business logic in separate functions.
│   │   │
│   │   ├── analyze/   # Analyze package containing the business logic and endpoints for analyzing an 
│   │   │              # application to be migrated, collecting the necessary data and outputting it 
│   │   │              # as an application profile.
│   │   │
│   │   ├── dockerize/ # Containerize the application to be migrated based on the application profile.
│   │
│   ├── utils/        # Contains helpful functions that can be reused throughout the code.

```

### CLI

```sh
cli/
│
├── cmd/             # Contains the commands for the CLI application
│   │
│   ├── analyze/     # Command for analyzing an application, residing in a VM, to collect the application data and construct an application profile with it
│   │
│   ├── dockerize/   # COmmand for containerization of an application, residing in a VM, based on the application profile, created by the analysis on the target VM
│
├── pkg/             # Contains reusable packages for the project
│   │
│   ├── utils/       # Utility functions shared across the project
```

## Software architecture

- A custom based architecture style, combining the following well-known and modern styles: N-layer, Modular monolith, Headless, and the Golang Package-based approach to separate the code base.

## Tech stack

- 🌐 **API**: Built with **Golang** using the **Gin** web framework for handling HTTP requests.
- 💻 **CLI**: Implemented in **Golang** using the **Cobra** library for building command-line tools.
- 🐧 **OS**: Runs on **Ubuntu**, a popular Linux distribution.
- ⚙️ **Scripting & Experiments**: Developed using **Ansible** for automation and **Golang** for prototyping and experimentation.

## Research stack

- 📄 **Research Paper**: Written in **LaTeX**
- 📚 **Citations**: Managed with **Scribbr** and **BibTeX**
- 📝 **Notes**: Organized using **Zettlr** and **Markdown**

## Run the system

### Analyze the target VM

The command analyzes the target VM, containing the application to be migrated. Based on the analysis it creates an application profile with the collected application files, ports, and services.

Example usage of the command:

```sh
go run main.go analyze \      
  --type=fs \                      
  --user=<username> \
  --host=<IP> \
  --privateKeyPath=<Path to the SSH key>
```

### Dockcerize the application

The command containerize the application using the application profile from the Analyze command. The result is a running Docker container with the application which is working and can be tested.

Example usage of the command:

```sh
go run main.go dockerize \      
  --dockerImageName=dockerized-vm \
  --dockerContainerName=dockerized-vm-container
```

## Research

The project was research-oriented and one of the core outcomes of it was delivering an academic [Research paper](./research-paper/Research-paper.pdf).

## Demo

The practical outcome of the project can be visualized through a [demo](./Demo.mp4) which outlines the core features.
