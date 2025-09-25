# phpbb-golang
A bulletin board inspired by phpBB. Written in Golang. It is using basic HTML, with optional JavaScript.

While the user interface is intentionally minimalistic, it is powered by a modern tech stack under the hood:
  - Golang: Fast, statically typed, and built for performance.
  - JSON-based data model.
  - Go Templates: Secure and efficient server-side rendering.
  - BBCode Support: Simple markup for user-friendly content formatting.
  - SQLite (default): Zero-config, file-based database for quick deployment setup. Can be migrated to PostgreSQL for production-grade performance.
  - Dockerized deployment (coming soon): Easy to deploy anywhere.
  - Visual Studio Code Dev Containers support: Seamless development environment — just open in Visual Studio Code and start coding.


# Screenshots

All text and content shown in the screenshots is for example purposes only.

![Main page](examples/myforum/screenshots/main.png?raw=true "Main page")
![Topics page](examples/myforum/screenshots/topics.png?raw=true "Topics page")
![Posts page](examples/myforum/screenshots/posts.png?raw=true "Posts page")


# Setting up Development Environment on Windows
 1. Setting up WSL2
      - Launch Command Prompt as Administrator.
      - Enable WSL with the latest Ubuntu version:
          ```
          wsl --install --distribution Ubuntu-24.04
          ```
      - Reboot.
      - Once rebooted, "wsl" command now provides much more options:
          ```
          wsl --help
          ```
      - Launch WSL:
          ```
          wsl
          ```
      - Upgrade all packages:
          ```
          $ apt list --installed
          $ sudo apt update
          $ sudo apt full-upgrade
          Do you want to continue? [Y/n] y
          $ apt list --installed
          ```

 2. Install Docker Daemon
      - Install Docker Daemon:
          ```
          $ sudo apt install -y docker.io
          ```
      - This automatically launches "/usr/bin/dockerd" in the background.
      - Add current user into "docker" group to be able to use docker command without sudo:
          ```
          $ id
          $ sudo usermod -aG docker $USER
          ```
      - Re-launch WSL.
      - Verify Docker installation:
          ```
          $ docker run hello-world
          Hello from Docker!
          This message shows that your installation appears to be working correctly.
          ```

 3. Develop Golang inside a Container:
      - Refs:
          - https://code.visualstudio.com/docs/remote/containers
          - https://github.com/microsoft/vscode-remote-try-go
      - Launch WSL:
          ```
          wsl
          ```
      - Launch Visual Studio Code:
          ```
          $ code . &
          ```
        This launches Visual Studio Code on Windows, and not Linux, but it's able to connect to /var/run/docker.sock. Magic?
      - Install "Dev Containers" extension in Visual Studio Code:
          - Click "View  >  Command Palette..."  (or press F1 or Ctrl+Shift+P).
          - Click "Extensions: Install Extensions".
          - Search Extensions in Marketplace: Dev Containers
          - Click "Install".
          - Restart the Visual Studio Code.
      - Reopen the code in Container:
          - Click "View  >  Command Palette...".
          - Click "Dev Containers: Reopen in Container".
          - Visual Studio Code shows "Dev Container: Go" (green background) at the bottom left of the screen.
      - Press F5 to run/debug the code.

    Troubleshooting:
      - To resolve issue "Docker returned an error. Make sure the Docker daemon is running and select an option how to proceed." and on terminal shows "Docker returned an error code ENOENT, message: Exectuable 'docker' not found on PATH":
          - Click "View  >  Command Palette...".
          - Click "Extensions: Show Enabled Extensions".
          - On "Dev Containers" extension, click "Manage (Gear icon)  >  Extension Settings".
          - On "Containers: Execute in WSL", make sure "[X] Controls whether CLI commands should always be executed in WSL." is checked.


# Development Tips
 1. To rebuild development environment:
      - Launch Visual Studio Code.
      - Click "View  >  Command Palette...".
      - Click "Dev Containers: Rebuild Container".
      - Run:
          ```
          go mod download
          ```

 2. To run tests:
      ```
      go test ./...
      ```

    To force running all tests:
      ```
      go clean -testcache
      ```

 3. To update "go.mod" and "go.sum" files:
      ```
      rm -f go.mod go.sum
      go mod init phpbb-golang
      go mod tidy
      ```
