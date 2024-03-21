--[[

  URL Behaviourlookup for rspamd
  
  Version: 1.0
  By Luc STEPNIEWSKI

--]]

local rspamd_logger = require "rspamd_logger"
local rspamd_http = require "rspamd_http"
local rspamd_url = require "rspamd_url"
local regexp = require "rspamd_regexp"

local N = "urlbehaviour"
local symbol_urlbehaviour = "URLBEHAVIOUR"
local symbol_urladdress = "URLBEHAVIOURADRESS"
local score_urlbehaviour = 15

local opts = rspamd_config:get_all_opt(N)
if not (opts and type(opts) == 'table') then
    rspamd_logger.infox(rspamd_config, 'Module is unconfigured')
    return
  end

-- Default settings
local api_url = "http://localhost:8088"
local cfg_timeout = 5

local function check_urlbehaviour(task)
    local function check_urlbehaviour_cb(err, code, body, headers)
        if err then
            rspamd_logger.errx('API call to URLBEHAVIOUR server failed: %s', err)
            return
        end

        rspamd_logger.infox('URLBEHAVIOUR returned body: %s', body)
        rspamd_logger.infox('URLBEHAVIOUR returned code: %s', code)
        rspamd_logger.infox('URLBEHAVIOUR returned headers: %s', headers)

        -- DO CODE !

        task:insert_result(symbol_urlbehaviour, 0, 'Got some URL behaviour header')
    end

    -- if not task:has_symbol(symbol_urlbehaviour) then
    --    return false
    --end

    -- Get the URLs from the body
    local URLs = task:get_urls({'https', 'http'}) or {}

    -- for _, url_iter in ipairs(URLs) do
    --    local function do_loop_iter()
    --       local url = url_iter
    --       rspamd_logger.infox('GOT URL=%s', url:tostring())
    --       -- check_url(url)
    --    end
    -- end

    local urlList = ''
    for _, url_iter in ipairs(URLs) do
        urlList = urlList .. '\n' .. tostring(url_iter)
    end

    if #URLs == 0 then
        rspamd_logger.infox('No url(s) in message (%s)', #URLs)
        return false
    end

    rspamd_logger.infox('querying URL BEHAVIOUR for urls %s', URLs)

    -- Getting some useful informations about the mail being checked
    local mail_from

    if task:has_from('smtp') then
        mail_from = task:get_from('smtp')[1]
    end
    if task:has_from('mime') then
        mail_from = task:get_from('mime')[1]
    end
    local client_host = task:get_from_ip()

    for _, uuu in ipairs(mail_from) do
        rspamd_logger.infox('mail_from content=%s, type is %s', uuu, type(uuu))
     end
 

    local content = tostring(client_host) .. '\n' .. tostring(mail_from) .. '\n'

    -- @TODO: need to make a JSON request, better to encode URLs?
    rspamd_http.request({
        -- url = api_url .. '?urls=' .. URLs,
        url = api_url,
        callback = check_urlbehaviour_cb,
        task = task,
        body = content .. urlList,
        -- headers={Header='Value', OtherHeader='Value'},
        mime_type = 'text/plain',
    })
end

-- MAIN

if opts then
    -- Loading options from config files
    if opts.url then
        api_url = opts.url
    end
    if opts.timeout then
        cfg_timeout = opts.timeout
    end

    rspamd_config:register_symbol({
        name = symbol_urlbehaviour,
        description = "Check URL Behaviour",
        score = score_urlbehaviour,
        callback = check_urlbehaviour
    })
    -- rspamd_config:register_dependency(symbol_urlbehaviour, symbol_urladdress)

    rspamd_logger.infox("%s module is configured and loaded (url=%s).", N, api_url)
else
    rspamd_logger.infox("%s module not configured.", N)
end