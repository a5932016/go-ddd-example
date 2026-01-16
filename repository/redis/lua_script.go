package redis

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

const (
	_LuaCMDSetManyHash = "set_many_hash"
)

func (m *MemRepository) execScript(cmd string, keys []string, args ...interface{}) (*redis.Cmd, error) {
	sha, err := m.getScriptSHA(cmd)
	if err != nil {
		return nil, err
	}

	return m.client.EvalSha(m.ctx, sha, keys, args...), nil
}

func (m *MemRepository) getScriptSHA(cmd string) (sha string, err error) {
	val, ok := m.scriptSHA1s.Load(cmd)
	if !ok {
		return m.loadScript(cmd)
	}

	return val.(string), nil
}

func (m *MemRepository) loadScript(cmd string) (sha string, err error) {
	filePath := fmt.Sprintf("repository/redis/lua_cmds/%s.lua", cmd)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	script := string(data)
	result := m.client.ScriptLoad(m.ctx, script)
	if err := result.Err(); err != nil {
		return "", err
	}

	sha = result.Val()
	m.scriptSHA1s.Store(_LuaCMDSetManyHash, sha)

	return
}
