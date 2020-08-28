Interval Restarter
=================

In order to restart a background ever-green program after every assigned period, you set an interval (as time duration) for it.

The following guide assumes your computer OS is a POSIX one.

How to install
-----------------

run this in a terminal

```shell
wget -qO- https://github.com/cxumol/interval-restarter/releases/download/v0.0.2/intrvl-r_0.0.2_linux_amd64.tar.gz | tar xzv && chmod +x intrvl-r
```

How to use
----------------

1. See `./cfg.yml.example` for details about configuration
2. Create your own `cfg.yml` or rename & edit from `./cfg.yml.example`
3. Run `./intrvl-r` in your terminal
4. Speicify a configuration file like `./intrvl-r -c "./my_config123.yml"`, if `cfg.yml` is used by other programs or you want more than one config files to use.

Why not crontab / cron job
-------------------

1. cron job works on the concept of calendar/clock, while this project is on the concept of timer;
2. A typical minimal time unit of cron job is minute, while that of this one is nanosecond;
3. cron job has a confusing working directory, while this one works under your current directory;
4. cron job has a troublesome logging way, while this one prints all the logs directly to your terminal

