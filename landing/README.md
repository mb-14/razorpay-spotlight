### README

```bash
python3 -m venv venv
. venv/bin/activate
pip install -r requirements.txt
FLASK_APP=app.py flask run
```

### Running on your VM

- Create the systemd service present in the systemd directory, do a daemon-reload and then enable the service and start it.
- Copy the nginx config over and follow the steps in the nginx directory, restart nginx then

