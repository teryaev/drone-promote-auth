load("@common//:steps.star", "notify_author")
load("@common//:utils.star", "ECR_URL", "retrieve_parameter")

def main(ctx):
    return [
        retrieve_parameter("DRONE_SLACK_BOT_TOKEN"),
        retrieve_parameter("DRONE_PEOPLEFORCE_API_KEY"),
        build_pipeline(ctx),
    ]

def build_pipeline(ctx):
    return {
        "kind": "pipeline",
        "name": "build and push drone ecs deploy image",
        "steps": [
            generate_tags_file(ctx),
            {
                "name": "build and push drone ecs deploy image",
                "image": "plugins/ecr",
                "settings": {
                    "registry": ECR_URL,
                    "repo": "drone-plugin/promote-auth",
                    "dockerfile": "Dockerfile",
                    "custom_dns": "169.254.169.253",
                },
            },
            notify_author(
                {"from_secret": "drone_slack_bot_token"},
                {"from_secret": "drone_peopleforce_api_key"},
            ),
        ],
        "trigger": {
            "branch": ["master"],
            "event": ["push"],
        },
    }

def generate_tags_file(ctx):
    commit_sha = ctx.build.commit[:6]

    return {
        "name": "generate tags file",
        "image": "alpine:3.11.5",
        "commands": [
            'echo -n "$(cat version),$DRONE_BUILD_NUMBER,latest,{}" > .tags'.format(commit_sha),
        ],
    }
