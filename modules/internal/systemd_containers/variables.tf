variable "containers" {
  description = "List of containers to start"
  type = list(object({
    # name becomes the container name and the service unit name prefix (e.g. <name>.service)
    name = string
    # ports uses short syntax [HOST:]CONTAINER[/PROTOCOL]
    # see also https://docs.docker.com/engine/network/#published-ports 
    ports = optional(list(string), [])
    # volume bind mount syntax is HOST:CONTAINER[:OPTIONS]
    volumes = optional(list(string), [])
    # env list of VAR=value pairs
    env        = optional(list(string), [])
    entrypoint = optional(string)
    network    = optional(string)
    pid        = optional(string)

    image = string
    # optional command to run (otherwise run the image's default)
    command = optional(string)

    # systemd service option settings
    # see also https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html#Options
    # where order is unimportant a single map can be used, otherwise order entries in a list
    service_options = optional(list(map(string)), [])
  }))

  validation {
    condition = alltrue([
      # Docker naming rule is a subset of systemd's unit name pattern.
      # To allow container DNS resolution we also limit the name to 63 characters.
      for c in var.containers : can(regex("^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,62}$", c.name))
    ]) && length(distinct(var.containers[*].name)) == length(var.containers)

    error_message = "Each container name must be unique, a valid docker container name and no more than 63 characters long."
  }
}
