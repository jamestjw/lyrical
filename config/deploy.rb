# config valid for current version and patch releases of Capistrano
# lock "~> 3.12.1"

set :application, "lyrical-bot"
set :repo_url, "https://github.com/jamestjw/lyrical.git"
ssh_options = { forward_agent: true }
set :gopath, deploy_to
set :pid_file, deploy_to+'/pids/PIDFILE'
set :linked_dirs, %w(configs log db audio-cache)


# Default branch is :master
# ask :branch, `git rev-parse --abbrev-ref HEAD`.chomp

# Default deploy_to directory is /var/www/my_app_name
# set :deploy_to, "/var/www/my_app_name"

# Default value for :format is :airbrussh.
# set :format, :airbrussh

# You can configure the Airbrussh format using :format_options.
# These are the defaults.
# set :format_options, command_output: true, log_file: "log/capistrano.log", color: :auto, truncate: :auto

# Default value for :pty is false
# set :pty, true

# Default value for :linked_files is []
# append :linked_files, "config/database.yml"

# Default value for linked_dirs is []
# append :linked_dirs, "log", "tmp/pids", "tmp/cache", "tmp/sockets", "public/system"

# Default value for default_env is {}
# set :default_env, { path: "/opt/ruby/bin:$PATH" }

# Default value for local_user is ENV['USER']
# set :local_user, -> { `git config user.name`.chomp }

# Default value for keep_releases is 5
set :keep_releases, 5

# Uncomment the following to require manually verifying the host key before first deploy.
# set :ssh_options, verify_host_key: :secure

namespace :deploy do
  task :restart do 
    on roles(:app) do
      old_pid = capture(:ps, :aux, '|', :grep, "bin/#{fetch(:application)}", '|', :grep, '-v grep', '|', :awk, "'{print $2}'", '|', :tail, '-n1')
      if old_pid && !old_pid.empty?
        info "Found PID: #{old_pid}"
        execute :kill, '-s', 'SIGINT', old_pid
        sleep 3
      end

      log_path = "#{release_path}/log/discordbot.log"
      execute "touch #{log_path}"
      execute "cd #{release_path}; screen -dmS discordbot sh -c './bin/lyrical-bot'"
    end
  end
end

after 'deploy:updated', 'go:build' do
  on roles(:app) do
    execute "cd #{release_path} && /usr/local/go/bin/go build -o #{release_path}/bin/lyrical-bot ."
  end
end
  
after 'go:build', 'go:stop-previous' do
  on roles(:app) do
    old_pid = capture(:ps, :aux, '|', :grep, "bin/#{fetch(:application)}", '|', :grep, '-v grep', '|', :awk, "'{print $2}'", '|', :tail, '-n1')
    if old_pid && !old_pid.empty?
      info "Found PID: #{old_pid}"
      execute :kill, '-s', 'SIGINT', old_pid
      sleep 3
    end
  end
end
  
after 'go:stop-previous', 'go:deploy-new' do
  on roles(:app) do
    log_path = "#{release_path}/log/discordbot.log"
    execute "touch #{log_path}"
    execute "cd #{release_path}; screen -dmS discordbot sh -c './bin/lyrical-bot'"
  end
end

