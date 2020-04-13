import threading,time,os,json
import win32gui,win32con,win32api
import requests,socket
import socketio as sio

init_start_num = 0
bootinfo_path = "Data/ClientBootInfo.txt"
Superviser_Log = "Data/SupervisorLog.txt"
file_c_time = 0
interval_sec = 3
window_hwnd = []
crash_time = {}
send_notify = []
crash_time_list = []

all_dead = True
has_new_request = False
start_time_stamp = int(time.time())

try:
    sio = sio.Client()
    sio.connect("//socket.path",socketio_path="/api/socket")
    my_ip = socket.gethostbyname(socket.gethostname())
    sio.emit("setIp",my_ip)
except Exception:
    print("Can't establish socket connetcion, this program will quit, error\n")
    quit()

def init_supervisor():
    global file_c_time
    global init_start_num
    global window_hwnd
    global all_dead

    all_dead = True
    current_file_m_time = os.path.getmtime(bootinfo_path)
    if(current_file_m_time > file_c_time):
        window_hwnd = []
        file_c_time = current_file_m_time

        f = open(bootinfo_path,"r")
        for line in f:
            if(line.find("Start Client")!=-1):
                init_start_num = int(line.split(':')[1])
            if(line.find("ClientInfo:")!=-1):
                window_hwnd.append([line.split(":")[1].split(",")[0],line.split(":")[1].split(",")[1][:-1],False])
        f.close()

        crash_time_list = []
        return True
    else:
        return False

    
def getPidByName():
    global all_dead
    global crash_time

    init_supervisor()
    result_list = []
    m, s = divmod(int(time.time()) - start_time_stamp, 60)
    h, m = divmod(m, 60)
    time_interval_str = "%d:%02d:%02d" % (h, m, s)
    time_str = time.strftime("%Y-%m-%d %H:%M:%S",time.localtime(int(time.time())))
    
    post_obj = {"data":[],"ip":my_ip}
    json_str = "Empty"

    for item in window_hwnd:
        client_info = item[:2]
        if(win32gui.GetWindowText(int(item[0])) == item[1]):
            all_dead = False
            
            client_info.append(True)
            client_info.append(time_interval_str)
            result_list.append(client_info)
        else:
            client_info.append(False)
            if(item[-1] == False):
                item[-1] = True
                crash_time[client_info[0]] = time_str

            if(crash_time[client_info[0]]):
                client_info.append(crash_time[client_info[0]])
            result_list.append(client_info)

    json_str = json.dumps(result_list)
    post_obj = {"data":result_list,"ip":my_ip}
    try:
        sio.emit("sendMonitor", post_obj)
        print(json_str)
        f_log = open(Superviser_Log,"a+")
        f_log.write(json_str+"\n")
        f_log.close()
        sio.on("pingServer")
    except Exception as e:
        print("Server error! "+e.__str__()+"\n")
    
    sendNotify(result_list)

    global t2
    t2.cancel()
    newtimer()
    t2.start()

def sendNotify(in_result_list):
    crash_time_list = list(crash_time.items())
    for i in crash_time_list:
        if(i[0] in send_notify):
            pass
        else:
            for item in in_result_list:
                if(item[0] == i[0]):
                    wind_name = item[1]
                    close_time = item[3]

            data = {
            "msgtype": "markdown",
            "markdown": {
            "content": "但丁训练客户端<font color='warning'>关闭</font>，请相关同事注意。\
            \n> 窗口:<font color='comment'>"+ wind_name +" </font> \
            \n> 句柄:<font color='comment'>"+ i[0] +" </font> \
            \n> 时间:<font color='comment'>"+ close_time +" </font> \
            \n> 服务器ip:<font color='comment'>"+ my_ip +" </font>"
                }
            }   

            r = requests.post("//callback.path" , json=data)
            send_notify.append(i[0])


@sio.on("pingServer")
def ping(data):
    sio.emit("pong", "//socket.path")

def newtimer():
    global t2
    t2 = threading.Timer(interval_sec,getPidByName)

t2 = threading.Timer(interval_sec,getPidByName)
t2.start()
print("Timer start")
