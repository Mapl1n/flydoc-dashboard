package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	ragURL, _ := url.Parse("http://localhost:8081")
	pipeURL, _ := url.Parse("http://localhost:8082")
	kgURL, _ := url.Parse("http://localhost:8083")

	// Reverse proxy with path rewrite: /rag/api/xxx → /api/xxx
	ragProxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = ragURL.Scheme
			r.URL.Host = ragURL.Host
		},
	}
	pipeProxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = pipeURL.Scheme
			r.URL.Host = pipeURL.Host
		},
	}
	kgProxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = kgURL.Scheme
			r.URL.Host = kgURL.Host
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Proxy: /s1/xxx → :8081/xxx
		if strings.HasPrefix(path, "/s1/") || strings.HasPrefix(path, "/s1") {
			if path == "/s1" {
				http.Redirect(w, r, "/s1/", http.StatusFound)
				return
			}
			r.URL.Path = strings.TrimPrefix(path, "/s1")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			ragProxy.ServeHTTP(w, r)
			return
		}

		// Proxy: /s2/xxx → :8082/xxx
		if strings.HasPrefix(path, "/s2/") || strings.HasPrefix(path, "/s2") {
			if path == "/s2" {
				http.Redirect(w, r, "/s2/", http.StatusFound)
				return
			}
			r.URL.Path = strings.TrimPrefix(path, "/s2")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			pipeProxy.ServeHTTP(w, r)
			return
		}

		// Proxy: /s3/xxx → :8083/xxx
		if strings.HasPrefix(path, "/s3/") || strings.HasPrefix(path, "/s3") {
			if path == "/s3" {
				http.Redirect(w, r, "/s3/", http.StatusFound)
				return
			}
			r.URL.Path = strings.TrimPrefix(path, "/s3")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			kgProxy.ServeHTTP(w, r)
			return
		}

		// Dashboard homepage
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(dashboardHTML))
	})

	addr := ":8484"
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  📂 FlyDoc 智能档案平台 — 统一入口")
	fmt.Println("========================================")
	fmt.Println("  Dashboard: http://localhost" + addr)
	fmt.Println("  /s1 → 📚 RAG 问答引擎      :8081")
	fmt.Println("  /s2 → 🔄 文档处理流水线    :8082")
	fmt.Println("  /s3 → 🧠 知识图谱构建器    :8083")
	fmt.Println("========================================")
	fmt.Println()

	log.Fatal(http.ListenAndServe(addr, nil))
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>📂 FlyDoc 智能档案平台</title>
<style>
:root {
  --bg: #0a0e1a;
  --card: #111827;
  --card-hover: #1a2332;
  --border: #1e293b;
  --text: #e2e8f0;
  --muted: #64748b;
  --accent: #3b82f6;
  --green: #22c55e;
  --purple: #8b5cf6;
  --orange: #f59e0b;
}
* { margin: 0; padding: 0; box-sizing: border-box; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: var(--bg);
  color: var(--text);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px 80px;
}
.header { text-align: center; margin-bottom: 52px; }
.header h1 { font-size: 34px; margin-bottom: 10px; display: flex; align-items: center; gap: 14px; justify-content: center; }
.header p { color: var(--muted); font-size: 16px; max-width: 600px; line-height: 1.6; }
.cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
  gap: 22px;
  max-width: 1140px;
  width: 100%;
}
.card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 32px;
  transition: all 0.25s;
  cursor: pointer;
  text-decoration: none;
  color: var(--text);
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}
.card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 4px;
  opacity: 0;
  transition: opacity 0.25s;
}
.card:nth-child(1)::before { background: var(--accent); }
.card:nth-child(2)::before { background: var(--green); }
.card:nth-child(3)::before { background: var(--purple); }
.card:hover { border-color: var(--accent); transform: translateY(-4px); box-shadow: 0 12px 40px rgba(0,0,0,0.3); }
.card:hover::before { opacity: 1; }
.card-icon { font-size: 44px; margin-bottom: 20px; }
.card h2 { font-size: 21px; margin-bottom: 12px; letter-spacing: -0.3px; }
.card p { color: var(--muted); font-size: 14px; line-height: 1.7; flex: 1; }
.card .tags { display: flex; gap: 6px; margin-top: 18px; flex-wrap: wrap; }
.tag { padding: 4px 10px; border-radius: 14px; font-size: 11px; font-weight: 500; }
.tag-blue { background: rgba(59,130,246,0.15); color: #93bbfd; }
.tag-green { background: rgba(34,197,94,0.15); color: #6ee7b7; }
.tag-purple { background: rgba(139,92,246,0.15); color: #c4b5fd; }
.tag-orange { background: rgba(245,158,11,0.15); color: #fcd34d; }
.card .action {
  margin-top: 22px;
  font-size: 14px;
  color: var(--accent);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}
.divider {
  width: 100%;
  max-width: 1140px;
  border-top: 1px solid var(--border);
  margin: 48px 0 24px;
  text-align: center;
  position: relative;
}
.divider span {
  position: relative;
  top: -10px;
  background: var(--bg);
  padding: 0 14px;
  color: var(--muted);
  font-size: 13px;
}
.link-row {
  max-width: 1140px;
  width: 100%;
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 14px;
}
.link-card {
  background: rgba(30,41,59,0.4);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  text-decoration: none;
  color: var(--text);
  transition: all 0.15s;
}
.link-card:hover { border-color: var(--accent); background: rgba(59,130,246,0.06); }
.link-card .name { font-size: 14px; font-weight: 500; }
.link-card .url { font-size: 11px; color: var(--muted); margin-top: 2px; }
.link-card .open-btn {
  padding: 6px 14px;
  border-radius: 6px;
  background: var(--accent);
  color: #fff;
  text-decoration: none;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}
@media (max-width: 768px) {
  .link-row { grid-template-columns: 1fr; }
  .header h1 { font-size: 26px; }
}
</style>
</head>
<body>

<div class="header">
  <h1>📂 FlyDoc 智能档案平台</h1>
  <p>基于 Go + ElasticSearch + Ollama + Tika 的智能档案处理解决方案<br>三大项目，一个入口，零依赖运行</p>
</div>

<div class="cards">

  <a class="card" href="/s1/" target="_blank">
    <div class="card-icon">📚</div>
    <h2>RAG 文档问答引擎</h2>
    <p>上传 PDF/DOCX → 智能解析 → 语义搜索 → AI 问答。<br>ES dense_vector 混合检索 + Ollama BGE-M3 本地 Embedding，零外部 API 费用。</p>
    <div class="tags">
      <span class="tag tag-blue">ES 混合检索</span>
      <span class="tag tag-green">Ollama Embedding</span>
      <span class="tag tag-purple">Tika 解析</span>
      <span class="tag tag-orange">LLM 流式问答</span>
    </div>
    <div class="action">🚀 打开项目 →</div>
  </a>

  <a class="card" href="/s2/" target="_blank">
    <div class="card-icon">🔄</div>
    <h2>文档处理流水线</h2>
    <p>Redis Stream 任务队列 + Go Worker Pool 并发处理。<br>WebSocket 实时推送进度，Pipeline 阶段编排，失败自动重试。</p>
    <div class="tags">
      <span class="tag tag-blue">Redis Stream</span>
      <span class="tag tag-green">Worker Pool</span>
      <span class="tag tag-purple">WebSocket</span>
      <span class="tag tag-orange">Tika 分类索引</span>
    </div>
    <div class="action">🚀 打开项目 →</div>
  </a>

  <a class="card" href="/s3/" target="_blank">
    <div class="card-icon">🧠</div>
    <h2>档案知识图谱构建器</h2>
    <p>纯 Go 中文 NER + 关系抽取 + ECharts 力导向图可视化。<br>自动识别人物/机构/合同/日期/金额，构建档案关联网络。</p>
    <div class="tags">
      <span class="tag tag-blue">中文 NER</span>
      <span class="tag tag-green">关系抽取</span>
      <span class="tag tag-purple">ECharts 力图</span>
      <span class="tag tag-orange">ES 图存储</span>
    </div>
    <div class="action">🚀 打开项目 →</div>
  </a>

</div>

<div class="divider"><span>📎 快捷访问</span></div>

<div class="link-row">
  <a class="link-card" href="/s1/" target="_blank">
    <div><div class="name">📚 RAG 问答引擎</div><div class="url">/s1 → :8081</div></div>
    <span class="open-btn">打开</span>
  </a>
  <a class="link-card" href="/s2/" target="_blank">
    <div><div class="name">🔄 文档处理流水线</div><div class="url">/s2 → :8082</div></div>
    <span class="open-btn">打开</span>
  </a>
  <a class="link-card" href="/s3/" target="_blank">
    <div><div class="name">🧠 知识图谱构建器</div><div class="url">/s3 → :8083</div></div>
    <span class="open-btn">打开</span>
  </a>
</div>

</body></html>`
