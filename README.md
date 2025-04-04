# gplan

gplan is an awesome CLI to-do list management application that helps you keep track of your projects and tasks efficiently.

## Table of Contents

1. [About The Project](#about-the-project)
   - [Built With](#built-with)
2. [Getting Started](#getting-started)
   - [Installation](#installation)
3. [Usage](#usage)
   - [Commands](#commands)
4. [Roadmap](#roadmap)
5. [Contributing](#contributing)
6. [License](#license)

## About The Project

gplan is a simple yet powerful CLI tool designed to help you manage your projects and tasks effectively from the terminal. Whether you're a developer, a project manager, or just someone who loves to keep things organized, gplan is the perfect tool for you.

### Built With

- [Go](https://golang.org/)
- [Cobra](https://github.com/spf13/cobra) - For building powerful modern CLI applications
- [SQLite](https://www.sqlite.org/index.html)
- [Tablewriter](https://github.com/olekukonko/tablewriter) - For table formatting in terminal
- [Treeprint](https://github.com/xlab/treeprint) - For tree formating in terminal

## Getting Started

To get a local copy up and running, follow these simple steps.

### Installation

You have several options to install gplan:

1. Download the official binary:

   - Get the appropriate binary for your operating system and computer architecture from the [releases page](https://codeberg.org/dynnian/gplan/releases).
   - Move the binary to a location in your PATH.

2. Install via Go:

   ```sh
   go install github.com/d4r1us-drk/gplan@latest
   ```

3. Install using Make:

   ```sh
   git clone https://github.com/d4r1us-drk/gplan.git
   cd gplan
   make install
   ```

## Usage

gplan allows you to manage projects and tasks with various commands. Below are some usage examples.

### Commands

- Create a new project:

  ```sh
  gplan new project -n "New Project" -d "Project Description"
  ```

- Create a new task with priority:

  ```sh
  gplan new task -n "New Task" -d "Task Description" -D "2024-08-15 23:00" -p "Existing Project" -P 1
  ```

  Priority levels: 1 (High), 2 (Medium), 3 (Low), 4 (None)

- Edit an existing project:

  ```sh
  gplan edit project 1 -n "Updated Project Name" -d "Updated Description"
  ```

- Edit a task's priority:

  ```sh
  gplan edit task 1 -P 2
  ```

- List all projects:

  ```sh
  gplan list projects
  ```

- List tasks by project:

  ```sh
  gplan list tasks -p "Project Name"
  ```

- Remove a project:

  ```sh
  gplan remove project 1
  ```

- Toggle task completion:

  ```sh
  gplan toggle 1
  ```

For detailed help, use the help command:

```sh
gplan help
```

## Roadmap

- [x] Add task and project management
- [x] Add priority levels for tasks
- [x] Implement Cobra framework for improved CLI structure
- [x] Add shell completion support
- [X] Add sub-tasks and sub-projects
- [X] Add a JSON output option to facilitate scripting
- [X] Use MVC Architecture and dependency injection
- [ ] Add a TUI interface
- [ ] Add a config file with customizable options, like database path, date-time format, etc.
- [ ] Add reminders and notifications (this would require a daemon)

See the [open issues](https://codeberg.org/dynnian/gplan/issues) for a full list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the GPLv3 License. See `LICENSE` for more information.
