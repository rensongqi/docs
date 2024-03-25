
# Index相关操作

```yaml
# Click the Variables button, above, to create your own variables.
GET ${exampleVariable1} // _search
{
  "query": {
    "${exampleVariable2}": {} // match_all
  }
}

# 创建一个指定类型的index
PUT test2
{
  "mappings": {
    "properties": {
      "name": {
        "type": "text"
      },
      "age": {
        "type": "long"
      },
      "addr": {
        "type": "text"
      }
    }
  }
}

# 创建index
# <index_name>/_doc/<doc_key>
PUT test2/_doc/1
{
  "name": "rsq",
  "age": 12,
  "addr": "cdscd"
}

# 获取指定document值
GET test2/_doc/1

# 通过_cat可以查看es当前的很多信息
GET _cat/health
GET _cat/indices

# 修改一个document的值
# method1: 覆盖原来的值
POST test2/_doc/1
{
  "name": "rsq1111",
  "age": 11,
  "addr": "axc"
}

POST test2/_doc/2
{
  "name": "rsq1111",
  "age": 22,
  "addr": "awx"
}

# method2: 更新指定的值
POST test2/_update/1
{
  "doc": {
      "name": "rsq111"
  }
}

# 删除索引
DELETE test3

# 删除dockument
DELETE test2/_doc/2

# 搜索指定key，精确匹配
GET test2/_search?q=name:rsq111

# 使用对象查询
# "_source"过滤指定的字段
GET test2/_search
{
  "query": {
    "match": {
      "name": "rsq1111"
    }
  },
  "_source": ["age", "addr"]
}

# 排序
# "order": "asc" 升序， desc 降序
GET test2/_search
{
  "query": {
    "match": {
      "name": "rsq1111"
    }
  },
  "sort": [
    {
      "age": {
        "order": "asc"
      }
    }
  ]
}

# 分页查询
# "from": 0,
# "size": 1
GET test2/_search
{
  "query": {
    "match": {
      "name": "rsq1111"
    }
  },
  "sort": [
    {
      "age": {
        "order": "asc"
      }
    }
  ],
  "from": 0,
  "size": 1
}

# bool查询
# must: and
# should: or
# not: !
# must_not: !
GET test2/_search
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "name": "rsq1111"
          }
        },
        {
          "match": {
            "age": "22"
          }
        }
      ]
    }
  }
}


# 查询结果过滤
# lt: 小于, lte: 小于等于, gt: 大于, gte: 大于等于, 
GET test2/_search
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "name": "rsq1111"
          }
        }
      ],
      "filter": [
        {
          "range": {
            "age": {
              "gt": 20
            }
          }
        }
      ]
    }
  }
}

# 构造新的数据

POST test3/_doc/1
{
  "name": "rsq1111",
  "age": 22,
  "addr": "awx",
  "hobby": ["旅游", "睡觉", "吃饭"]
}

POST test3/_doc/2
{
  "name": "zhangsan1111",
  "age": 11,
  "addr": "awc",
  "hobby": ["睡觉", "吃饭"]
}

POST test3/_doc/3
{
  "name": "lisi1111",
  "age": 121,
  "addr": "awxc",
  "hobby": ["吃饭", "唱歌"]
}

# 支持在数组中查找，多个条件以空格为分隔符
GET test3/_search
{
  "query": {
    "match": {
      "hobby": "旅 唱"
    }
  }
}

# term查询时是通过倒排索引进行精确查询
GET test3/_search
{
  "query": {
    "term": {
      "age": 121
    }
  }
}

# text类型会被分词器解析，keyword类型是不会被分词器解析的
GET _analyze
{
  "analyzer": "standard",
  "text": "我爱学习 golang"
}

GET _analyze
{
  "analyzer": "keyword",
  "text": "我爱学习 golang"
}


# 查询高亮
GET test3/_search
{
  "query": {
    "match": {
      "name": "lisi1111"
    }
  },
  "highlight": {
    "pre_tags": "<p class='key' style'color:red'>",
    "post_tags": "</p>",
    "fields": {
      "name": {}
    }
  }
}
```