package main

import (
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// OrganizationsFormToOrganizationsV2 함수는 form 문자를 받아서 []Organization 을 생성한다.
func OrganizationsFormToOrganizationsV2(client *mongo.Client, s string) ([]Organization, error) {
	var results []Organization
	orgs := strings.Split(s, ":")
	for _, org := range orgs {
		parts := strings.Split(org, ",")
		if len(parts) != 6 { // [ Primary여부(true||false), Division, Department, Team, Role, Position ] 총 6개의 마디로 되어있다.
			continue
		}
		org := Organization{}
		if parts[0] == "true" {
			org.Primary = true
		} else {
			org.Primary = false
		}
		if parts[1] != "unknown" {
			division, err := getDivisionV2(client, parts[1])
			if err != nil {
				return results, err
			}
			org.Division = division
		}
		if parts[2] != "unknown" {
			department, err := getDepartmentV2(client, parts[2])
			if err != nil {
				return results, err
			}
			org.Department = department
		}
		if parts[3] != "unknown" {
			team, err := getTeamV2(client, parts[3])
			if err != nil {
				return results, err
			}
			org.Team = team
		}
		if parts[4] != "unknown" {
			role, err := getRoleV2(client, parts[4])
			if err != nil {
				return results, err
			}
			org.Role = role
		}
		if parts[5] != "unknown" {
			position, err := getPositionV2(client, parts[5])
			if err != nil {
				return results, err
			}
			org.Position = position
		}
		results = append(results, org)
	}
	return results, nil
}
