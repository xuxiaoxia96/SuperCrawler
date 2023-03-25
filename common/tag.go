package common

import "strings"

func AddExamTag(title string) string{
	tags := []string{"公务员", "事业单位", "教师", "医疗", "选调", "遴选", "三支一扶", "大学生村官", "基层工作者", "银行", "国企", "公益性岗位"}
	for _,tag := range tags{
		if strings.Contains(title, tag){
			return tag
		}
	}

	return ""
}
