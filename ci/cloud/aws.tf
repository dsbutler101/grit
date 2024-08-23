resource "aws_iam_user" "grit-tester" {
  name = "grit-tester"
}

resource "aws_iam_access_key" "grit-tester" {
  user = aws_iam_user.grit-tester.name
}

resource "aws_iam_user_policy_attachment" "grit-tester-administrator" {
  # kics-scan ignore-line "Grandfathered in, might rework in the future. Ref: https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/merge_requests/103#note_2037121940"
  user       = aws_iam_user.grit-tester.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}
