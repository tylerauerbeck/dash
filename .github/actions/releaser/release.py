import datetime, croniter, json, os, git
from crontab import CronTab

# Grab Github Event Path and Workspace to get cron schedule and current repo
#EVENT = os.getenv('GITHUB_EVENT_PATH')
WORKSPACE = os.getenv('GITHUB_WORKSPACE')

# Read in the event so we can parse the cron schedule
#with open(EVENT) as ev:
#  data = json.load(ev)
#  cron_sched = data['schedule']

cron_sched = '0 12 * * 5'

# Create an empty crontab
cron = CronTab()

# Create a dummy job -- we just need to be able to evaluate a cron schedule
job = cron.new(command='/usr/bin/echo')
job.setall(cron_sched)

# Create a schedule so that we can look at the previous run
schedule = job.schedule(date_from=datetime.datetime.now())
last_run = schedule.get_prev()

print("Checking to see if there have been commits since " + str(last_run))

# Hook into current repo and look at the git logs
repo = git.Repo(WORKSPACE)
logs = repo.git.log('--oneline', '--pretty=format:" %h %s by %an"', '--no-merges', '--since=' + str(last_run)).split('\n')

# If there have been commits in the repo since the last release
if len(logs) >= 1:
  print("There have been commits since the last release. Creating a new release tag.")

# Grab or calculate the current release tag
  if len(repo.tags) == 0:
    latest_release = "v0.0.0"
  else:
    latest_release = repo.tags[-1]

# Increment the patch release by 1 and create the new release tag
  new_release = '.'.join(str(latest_release).split('.')[:-1]) + '.' + str(int(str(latest_release).split('.')[-1]) + 1)
#  new_tag = repo.create_tag(new_release, message='release "{0}"'.format(new_release))
  print(new_release)

# Run the release process
#  os.system('cd $GITHUB_WORKSPACE && make release')
  os.system('echo ::set-output name=version::' + new_release)
else:
  print("There have been no commits since the last release")
