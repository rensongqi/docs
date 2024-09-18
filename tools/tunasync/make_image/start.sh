#!/bin/sh
tunasync manager --config /data/tunasync/manager.conf &
tunasync worker --config /data/tunasync/resources.conf
