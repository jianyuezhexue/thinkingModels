import subprocess
import time
import sys

# 启动前端服务
print("启动前端服务...")
frontend_proc = subprocess.Popen(
    ["pnpm", "dev:ele"],
    cwd="/Users/wang/code/thinkingModels/frontend",
    stdout=subprocess.PIPE,
    stderr=subprocess.STDOUT,
    text=True
)

# 等待服务启动
print("等待服务启动...")
time.sleep(15)

# 检查服务输出
print("\n=== 前端服务日志 ===")
try:
    # 非阻塞读取
    import select
    if select.select([frontend_proc.stdout], [], [], 0)[0]:
        output = frontend_proc.stdout.read()
        print(output)
except:
    # 回退到直接读取
    pass

# 检查进程状态
if frontend_proc.poll() is None:
    print("\n前端服务正在运行 (PID: {})".format(frontend_proc.pid))
else:
    print("\n前端服务启动失败，退出码: {}".format(frontend_proc.returncode))

# 保持运行一段时间以便检查
print("\n服务已启动，按 Ctrl+C 停止...")
try:
    time.sleep(30)
except KeyboardInterrupt:
    pass
finally:
    frontend_proc.terminate()
    print("服务已停止")