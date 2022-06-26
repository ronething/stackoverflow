# coding=utf-8

import time
from loguru import logger

if __name__ == '__main__':
    logger.add("main.log") # 测试是否可以立刻生成对应文件
    logger.info("hello world")
    logger.info("main func run")
    time.sleep(10)
