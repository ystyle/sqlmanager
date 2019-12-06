### GetStudentByID2

```sql
	select * from student where id = {{.}}
```

### GetStudentByID
>通过ID查找学生
```sql
	select * from student where id = {{.id}}
```

### GetStudentByID3

>通过ID查找学生
```sql
	select * from student where id = {{.id}}
```

### QueryStudentMessageTpl

>查询学生消息

```sql
select 
       victoria_user_message.type,
       victoria_user_message.title,
       victoria_user_message.title_en,
       victoria_user_message.content,
       victoria_user_message.content_en,
       victoria_user_message.is_reply,
       victoria_user_message.range,
       victoria_user_message.operation,
       victoria_user_message.operation_detail,
       victoria_user_message.expire_time,
       victoria_user_message.attachments,
       victoria_user_message_logs.id,
       victoria_user_message_logs.publish_time,
       victoria_user_message_logs.message_id,
       victoria_user_message_logs.student_id,
       victoria_user_message_logs.readed_time,
       victoria_user_message_logs.Replied
from victoria_user_message_logs
         left join victoria_user_message on victoria_user_message_logs.message_id = victoria_user_message.id
where victoria_user_message_logs.deleted_at is null and  victoria_user_message.deleted_at is null
{{if .ID}}
and victoria_user_message_logs.id = {{ .ID }} 
{{end}}
{{if .StudentID}}
and victoria_user_message_logs.student_id = '{{ .StudentID }}'
{{end}}
{{if .Type}}
and victoria_user_message.type = {{ .Type }}
{{end}}
{{if .Keyword}}
and (
 victoria_user_message.title like '%{{.Keyword}}%'
 or victoria_user_message.title_en like '%{{.Keyword}}%'
 or victoria_user_message.content like '%{{.Keyword}}%'
 or victoria_user_message.content_en like '%{{.Keyword}}%'
)
{{end}}
{{if eq .ReadStatus 1}}
victoria_user_message_logs.readed_time is not null
{{else if eq .ReadStatus 2}}
victoria_user_message_logs.readed_time is null
{{end}}
```