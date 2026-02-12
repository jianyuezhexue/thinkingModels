from PIL import Image, ImageDraw, ImageFont
import os

def create_title_slide_background():
    """创建标题幻灯片背景"""
    width, height = 1920, 1080
    img = Image.new('RGB', (width, height), color=(245, 245, 245))  # 浅灰色背景
    draw = ImageDraw.Draw(img)
    
    # 创建渐变效果
    for y in range(height):
        r = int(245 - (y / height) * 20)
        g = int(245 - (y / height) * 20)
        b = int(255 - (y / height) * 30)
        draw.line([(0, y), (width, y)], fill=(r, g, b))
    
    # 添加装饰元素
    draw.ellipse([50, 50, 300, 300], outline=(100, 150, 200), width=2)
    draw.ellipse([1620, 100, 1870, 350], outline=(200, 100, 150), width=2)
    
    img.save('/Users/wang/code/thinkingModels/title_bg.png')
    print("标题背景图片已创建")

def create_content_background():
    """创建内容幻灯片背景"""
    width, height = 1920, 1080
    img = Image.new('RGB', (width, height), color=(255, 255, 255))
    draw = ImageDraw.Draw(img)
    
    # 左侧边栏
    draw.rectangle([0, 0, 300, height], fill=(240, 240, 240))
    
    # 底部装饰条
    draw.rectangle([0, 980, width, height], fill=(230, 240, 250))
    
    # 添加几何图形装饰
    for i in range(5):
        x = 50 + i * 100
        y = 800 + (i % 2) * 50
        draw.ellipse([x, y, x+50, y+50], fill=(200, 220, 240))
    
    img.save('/Users/wang/code/thinkingModels/content_bg.png')
    print("内容背景图片已创建")

def create_chart_image():
    """创建简单的图表图像"""
    width, height = 800, 600
    img = Image.new('RGB', (width, height), color=(255, 255, 255))
    draw = ImageDraw.Draw(img)
    
    # 绘制标题
    try:
        font = ImageFont.truetype("Arial.ttf", 24)
    except:
        font = ImageFont.load_default()
    
    draw.text((300, 30), "商业模式概览", fill=(0, 51, 102), font=font)
    
    # 绘制简单的柱状图
    bar_width = 80
    bar_spacing = 40
    bars = [
        ("免费用户", 10),
        ("年度会员", 89),
        ("终生会员", 399),
        ("企业席位", 20)
    ]
    
    max_value = max(bar[1] for bar in bars)
    chart_height = 400
    chart_width = 600
    chart_x = 100
    chart_y = 150
    
    # X轴
    draw.line([(chart_x, chart_y + chart_height), (chart_x + chart_width, chart_y + chart_height)], fill=(0, 0, 0), width=2)
    
    # Y轴
    draw.line([(chart_x, chart_y), (chart_x, chart_y + chart_height)], fill=(0, 0, 0), width=2)
    
    # 绘制柱子
    for i, (label, value) in enumerate(bars):
        bar_x = chart_x + i * (bar_width + bar_spacing) + bar_spacing
        bar_height = (value / max_value) * chart_height * 0.8
        bar_y = chart_y + chart_height - bar_height
        
        # 绘制柱子
        draw.rectangle([bar_x, bar_y, bar_x + bar_width, chart_y + chart_height], fill=(70, 130, 180))
        
        # 绘制标签
        draw.text((bar_x, chart_y + chart_height + 10), label, fill=(0, 0, 0), font=font)
        
        # 绘制数值
        draw.text((bar_x + 10, bar_y - 25), f"¥{value}", fill=(0, 0, 0), font=font)
    
    img.save('/Users/wang/code/thinkingModels/chart.png')
    print("图表图片已创建")

if __name__ == "__main__":
    create_title_slide_background()
    create_content_background()
    create_chart_image()