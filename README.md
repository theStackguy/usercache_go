# usercache_go
memory cache management system designed to optimize user data retrieval and storage. 

# 🧠 User Memory Cache Manager (C#)

A lightweight and efficient memory cache management system for handling user-specific data in C# applications. This project provides a robust caching layer to reduce database load, improve response times, and enhance user experience.

## 🚀 Features

- **User-centric caching**: Store and retrieve data per user session or identifier.
- **Configurable expiration**: Set custom TTL (Time-To-Live) for cache entries.
- **Thread-safe operations**: Safe for concurrent access in multi-threaded environments.
- **Pluggable cache backend**: Easily switch between in-memory and distributed cache providers.
- **Monitoring support**: Track cache hits, misses, and evictions.

## 🛠️ Usage

### 1. Add the Cache Manager

```go
using System;
using System.Runtime.Caching;

public class UserCacheManager
{
    private readonly MemoryCache _cache = MemoryCache.Default;

    public void Set(string userId, object data, TimeSpan expiration)
    {
        var policy = new CacheItemPolicy { AbsoluteExpiration = DateTimeOffset.Now.Add(expiration) };
        _cache.Set(userId, data, policy);
    }

    public object Get(string userId)
    {
        return _cache.Get(userId);
    }

    public void Remove(string userId)
    {
        _cache.Remove(userId);
    }
}

