export function loadArticlesApi(n) {
    console.log('loadArticles')
    const articles = []
    for (let i = 0; i < n; i++) {
        articles.push({
            index: i + 1,
            title: `文章${i + 1}`,
            progress: Math.floor(Math.random() * 100),
            status: Math.random() > 0.5 ? '成功' : '失败',
        })
    }
    articles.push({
        index: 99,
        title: `文章99`,
        progress: 100,
        status: Math.random() > 0.5 ? '成功' : '失败',
    })
    return articles
}

export function loadPlatformsApi(n) {
    console.log('loadPlatform')
    const platforms = []
    for (let i = 0; i < n; i++) {
        platforms.push({
                index: i + 1,
                name: `平台${i + 1}`,
                connected: Math.random() > 0.5,
                disabled: Math.random() > 0.5,
            }
        )
    }
    return platforms

}
