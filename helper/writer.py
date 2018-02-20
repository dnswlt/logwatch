import string
import time
from datetime import datetime
import random
import sys


chars = string.ascii_lowercase + string.digits

with open("output.log", "a") as f1, open("output2.log", "a") as f2:
    while True:
        f1, f2 = f2, f1
        t = datetime.now().strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]
        msg = ''.join(random.choices(chars, k=random.randint(80, 200)))
        level = "INFO" if random.random() > 0.05 else "ERROR"
        f1.write("%s %s %s\n" % (t, level, msg))
        f1.flush()
        time.sleep(0.01)
        