# 查看版本
version:
	@grep -p 'Version = ' utils.go|cut -f2 -d'"'

# 提交版本
tag:
	@awk -F '"' '/Version/ {print $$2;system("git tag "$$2);system("git push origin "$$2)}' utils.go

# 删除远端版本 慎用
dropTag:
	@awk -F '"' '/Version/ {print $$2;system("git tag -d "$$2);system("git push origin :refs/tags/"$$2)}' utils.go

