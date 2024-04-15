# example

Command example is the malatd service project.
<br>The framework reference: https://github.com/swxctx/malatd

## API Desc

### V1_Chat_Do

Do handler

- URI: `/example/v1/chat/do`
- METHOD: `POST`
- QUERY:
- BODY:

	```js
	{
		"platform": -0,	// {int} // 厂商[1: 百度 2: 讯飞 3: 智谱]
		"stream": false,	// {bool} // 是否流式
		"content": ""	// {string} // 对话内容
	}
	```

- RESPONSE:

	```js
	{
		"message": ""	// {string} 
	}
	```





<br>

*This is a project created by `malatd gen` command.*

*[About Malatd Command](https://github.com/swxctx/malatd)*

## Error List

|Code|Message|
|------|------|
|100001| "Invalid Parameter"|
