#!/bin/bash
set -euo pipefail

export PODMAN_COMPOSE_WARNING_LOGS=false
export MYSQL_PWD="root"

start_podman() {
	echo "Starting services..."
	podman compose start >/dev/null
}

wait_for_database() {
	echo "Waiting for database to accept connections..."
	until podman exec snippetbox_db mysqladmin ping --host=localhost --user=root --silent >/dev/null; do
		sleep 1
	done
}

start_modd() {
	modd
}

stop_podman() {
	echo "Stopping services..."
	podman compose stop >/dev/null
}

remove_binary() {
	echo "Removing application binary..."
	rm "$(pwd)/web"
}

cleanup() {
	stop_podman
	remove_binary
}

trap cleanup SIGINT SIGTERM

start_podman
wait_for_database
start_modd
