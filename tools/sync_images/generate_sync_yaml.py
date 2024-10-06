import os
import re
import yaml
import requests
import subprocess
import json
from distutils.version import LooseVersion
from urllib.parse import quote
import concurrent.futures
import datetime

# 基本配置
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
CONFIG_FILE = os.path.join(BASE_DIR, 'sync_images/config.yaml')
CUSTOM_SYNC_FILE = os.path.join(BASE_DIR, 'custom_sync.yaml')
LOG_FILE_PATH = "/var/log/sync_images"


def is_exclude_tag(tag):
    """
    排除tag
    :param tag:
    :return:
    """
    excludes = ['alpha', 'dev', 'test', 'ppc64le', 'arm', 's390x', 'beta','SNAPSHOT', 'debug', 'debian', 'sha256-']

    for e in excludes:
        if e.lower() in tag.lower():
            return True
        if str.isalpha(tag):
            if tag == 'latest':
                return False
            else:
                return True
        if len(tag) >= 40:
            return True

    pattern = re.compile(r'(.*-(full|minimal)$)|(-[a-z0-9]{9,}$)')
    if pattern.match(tag):
        return True

    return False


# 通过给定指定仓库信息，获取repo对应的tag
def fetch_harbor_repo_tags(harbor_url, project, repo_list, username, password, repo_page_size):
    repo_page = 1
    temp_tags_lists = []
    repo_tags = []
    repo_name = repo_list["name"].removeprefix(project+'/')
    encoded_repo_name = quote(quote(repo_name, safe=''), safe='')
    while True:
        repo_url = f"{harbor_url}/api/v2.0/projects/{project}/repositories/{encoded_repo_name}/artifacts?page={repo_page}&page_size={repo_page_size}"
        response = requests.get(repo_url, auth=(username, password))
        response.raise_for_status()
        page_data = response.json()
        if not page_data:
            break
        temp_tags_lists.extend(page_data)
        repo_page += 1
    # 提取标签
    for artifact in temp_tags_lists:
        if 'tags' in artifact and artifact['tags'] is not None:
            for tag in artifact['tags']:
                repo_tags.append(tag['name'])
    return repo_name, repo_tags


def get_harbor_tags(project):
    """
    :param project:
    :return:
    """
    print("[get_harbor_tags] project: ", project)
    # Harbor 服务器信息
    harbor_url = "https://harbor.rsq.cn"
    username = "harbor"
    password = "GphVjjGe0"
    # 分页获取数据
    project_page = 1
    project_page_size = 100
    repositories = []
    while True:
        # 获取指定项目中对应的所有仓库信息
        repo_url = f"{harbor_url}/api/v2.0/projects/{project}/repositories?page={project_page}&page_size={project_page_size}"
        response = requests.get(repo_url, auth=(username, password))
        response.raise_for_status()
        page_data = response.json()
        if not page_data:
            break
        repositories.extend(page_data)
        project_page += 1
    # 初始化字典来存储 repo_list 和 repo_tags 的对应关系
    repo_tags_map = {}
    # 使用线程池获取所有仓库的标签
    with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
        future_to_repo = {executor.submit(fetch_harbor_repo_tags, harbor_url, project, repo_list, username, password, project_page_size): repo_list for repo_list in repositories}
        for future in concurrent.futures.as_completed(future_to_repo):
            repo_list = future_to_repo[future]
            try:
                repo_name, repo_tags = future.result()
                repo_tags_map[repo_name] = repo_tags
                # print(f"Tags for repository {repo_name}: {repo_tags}")
            except Exception as exc:
                print(f"Repository {repo_list['name']} generated an exception: {exc}")
    return repo_tags_map


def get_repo_gcr_tags(harbor_image_tags, image, host="k8s.gcr.io"):
    """
    获取 gcr.io repo 最新的 tag
    :param host:
    :param image:
    :param limit:
    :return:
    """

    headers = {
        'User-Agent': 'docker/19.03.12 go/go1.13.10 git-commit/48a66213fe kernel/5.8.0-1.el7.elrepo.x86_64 os/linux arch/amd64 UpstreamClient(Docker-Client/19.03.12 \(linux\))'
    }

    tag_url = "https://{host}/v2/{image}/tags/list".format(host=host, image=image)

    tags = []
    tags_data = []
    manifest_data = []

    try:
        tag_rep = requests.get(url=tag_url, verify=False, timeout=5, headers=headers)
        tag_req_json = tag_rep.json()
        manifest_data = tag_req_json['manifest']
    except Exception as e:
        print('[Get tag Error]', e)
        return tags

    for manifest in manifest_data:
        sha256_data = manifest_data[manifest]
        sha256_tag = sha256_data.get('tag', [])
        if len(sha256_tag) > 0:
            # 排除 tag
            if is_exclude_tag(sha256_tag[0]):
                continue
            tags_data.append({
                'tag': sha256_tag[0],
                'timeUploadedMs': sha256_data.get('timeUploadedMs')
            })
    tags_sort_data = sorted(tags_data, key=lambda i: i['timeUploadedMs'], reverse=True)
    for t in tags_sort_data:
        # 去除同步过的
        if image in harbor_image_tags:
            if t['tag'] in harbor_image_tags[image]:
                continue
            tags.append(t['tag'])

    print(f'[{image} repo tag]', tags)
    return tags


def get_repo_quay_tags(harbor_image_tags, image):
    """
    获取 quay.io repo 最新的 tag
    :param image:
    :return:
    """

    headers = {
        'User-Agent': 'docker/19.03.12 go/go1.13.10 git-commit/48a66213fe kernel/5.8.0-1.el7.elrepo.x86_64 os/linux arch/amd64 UpstreamClient(Docker-Client/19.03.12 \(linux\))'
    }

    tag_url = "https://quay.io/api/v1/repository/{image}/tag/?onlyActiveTags=true&limit=100".format(image=image)

    tags = []
    tags_data = []
    try:
        tag_rep = requests.get(url=tag_url, verify=False, timeout=5, headers=headers)
        tag_req_json = tag_rep.json()
        manifest_data = tag_req_json['tags']
    except Exception as e:
        print('[Get tag Error]', e)
        return tags

    for manifest in manifest_data:
        name = manifest.get('name', '')
        # 排除 tag
        if is_exclude_tag(name):
            continue

        tags_data.append({
            'tag': name,
            'start_ts': manifest.get('start_ts')
        })

    tags_sort_data = sorted(tags_data, key=lambda i: i['start_ts'], reverse=True)

    for t in tags_sort_data:
        # 去除同步过的
        if image in harbor_image_tags:
            if t['tag'] in harbor_image_tags[image]:
                continue
            tags.append(t['tag'])

    print('[repo tag]', tags)
    return tags


def get_repo_elastic_tags(image, limit=20):
    """
    获取 elastic.io repo 最新的 tag
    :param image:
    :param limit:
    :return:
    """
    token_url = "https://docker-auth.elastic.co/auth?service=token-service&scope=repository:{image}:pull".format(
        image=image)
    tag_url = "https://docker.elastic.co/v2/{image}/tags/list".format(image=image)

    tags = []
    tags_data = []
    manifest_data = []

    headers = {
        'User-Agent': 'docker/19.03.12 go/go1.13.10 git-commit/48a66213fe kernel/5.8.0-1.el7.elrepo.x86_64 os/linux arch/amd64 UpstreamClient(Docker-Client/19.03.12 \(linux\))'
    }

    try:
        token_res = requests.get(url=token_url, verify=False, timeout=5, headers=headers)
        token_data = token_res.json()
        access_token = token_data['token']
    except Exception as e:
        print('[Get repo token]', e)
        return tags

    headers['Authorization'] = 'Bearer ' + access_token

    try:
        tag_rep = requests.get(url=tag_url, verify=False, timeout=5, headers=headers)
        tag_req_json = tag_rep.json()
        manifest_data = tag_req_json['tags']
    except Exception as e:
        print('[Get tag Error]', e)
        return tags

    for tag in manifest_data:
        # 排除 tag
        if is_exclude_tag(tag):
            continue
        tags_data.append(tag)

    tags_sort_data = sorted(tags_data, key=LooseVersion, reverse=True)

    # limit tag
    tags_limit_data = tags_sort_data[:limit]

    harbor_image_tags = get_harbor_tags("docker.elastic.co")
    for t in tags_limit_data:
        # 去除同步过的
        if t in harbor_image_tags:
            continue

        tags.append(t)

    print('[repo tag]', tags)
    return tags


def get_repo_ghcr_tags(image, limit=20):
    """
    获取 ghcr.io repo 最新的 tag
    :param image:
    :param limit:
    :return:
    """
    token_url = "https://ghcr.io/token?service=ghcr.io&scope=repository:{image}:pull".format(
        image=image)

    tag_url = "https://ghcr.io/v2/{image}/tags/list".format(image=image)

    tags = []
    tags_data = []

    headers = {
        'User-Agent': 'docker/19.03.12 go/go1.13.10 git-commit/48a66213fe kernel/5.8.0-1.el7.elrepo.x86_64 os/linux arch/amd64 UpstreamClient(Docker-Client/19.03.12 \(linux\))'
    }

    try:
        token_res = requests.get(url=token_url, verify=False, timeout=5, headers=headers)
        token_data = token_res.json()
        print("token_data", token_url, token_data)
        access_token = token_data['token']
    except Exception as e:
        print('[Get repo token]', e)
        return tags

    headers['Authorization'] = 'Bearer ' + access_token

    try:
        tag_rep = requests.get(url=tag_url, verify=False, timeout=5, headers=headers)
        tag_req_json = tag_rep.json()
        manifest_data = tag_req_json['tags']
    except Exception as e:
        print('[Get tag Error]', e)
        return tags

    for tag in manifest_data:
        # 排除 tag
        if is_exclude_tag(tag):
            continue
        tags_data.append(tag)

    tags_sort_data = sorted(tags_data, key=LooseVersion, reverse=True)

    # limit tag
    tags_limit_data = tags_sort_data[:limit]

    harbor_image_tags = get_harbor_tags("ghcr.io")
    for t in tags_limit_data:
        # 去除同步过的
        if t in harbor_image_tags:
            continue

        tags.append(t)

    print('[repo tag]', tags)
    return tags


def get_image_tags_skopeo(repository):
    command = ["skopeo", "list-tags", "--no-creds", f"docker://{repository}"]
    try:
        # 调用 skopeo 命令
        result = subprocess.run(command, capture_output=True, text=True, check=True)
        output = result.stdout
        # 解析 JSON 输出
        data = json.loads(output)
        return data.get("Tags", [])
    except subprocess.CalledProcessError as e:
        print(f"Error occurred: {e}")
        return []
    except json.JSONDecodeError as e:
        print(f"Failed to parse JSON: {e}")
        return []


def get_nvcr_io_tags(harbor_image_tags, image):
    tags = []
    tags_data = []
    manifest_data = get_image_tags_skopeo('ngc.nju.edu.cn/'+image)

    for tag in manifest_data:
        if is_exclude_tag(tag):
            continue
        tags_data.append(tag)

    for t in tags_data:
        # 去除同步过的
        if image in harbor_image_tags:
            if t in harbor_image_tags[image]:
                continue
            tags.append(t)
        else:
            tags = tags_data
    return tags


def get_docker_io_tags(harbor_image_tags, image):
    # username = ""
    # image_name = ""
    # headers = {
    #     'User-Agent':
    #     'docker/19.03.12 go/go1.13.10 git-commit/48a66213fe kernel/5.8.0-1.el7.elrepo.x86_64 os/linux arch/amd64 UpstreamClient(Docker-Client/19.03.12 \(linux\))'
    # }
    # namespace_image = image.split('/')
    # if len(namespace_image) == 1:
    #     username = 'library'
    #     image_name = namespace_image[0]
    # elif len(namespace_image) > 1:
    #     username = namespace_image[0]
    #     image_name = namespace_image[1]
    #
    # print("namespace_image: ", namespace_image)
    # tag_url = "https://hub.docker.com/v2/namespaces/{user}/repositories/{image}/tags".format(
    #     user=username, image=image_name)
    # print(tag_url)

    tags = []
    tags_data = []

    # try:
    #     tag_rep = requests.get(url=tag_url, verify=False, timeout=5, headers=headers)
    #     tag_req_json = tag_rep.json()
    #     manifest_data = tag_req_json['results']
    # except Exception as e:
    #     print('[Get tag Error]', e)
    #     return tags

    manifest_data = get_image_tags_skopeo('docker.io/'+image)

    for tag in manifest_data:
        if is_exclude_tag(tag):
            continue
        tags_data.append(tag)
    for t in tags_data:
        # 去除同步过的
        if image in harbor_image_tags:
            if t in harbor_image_tags[image]:
                continue
            tags.append(t)
        else:
            tags = tags_data
    return tags


def get_repo_tags(repo, image, harbor_tags):
    """
    获取 repo 最新的 tag
    :param repo:
    :param image:
    :return:
    """
    tags_data = []
    if repo == 'gcr.io':
        tags_data = get_repo_gcr_tags(harbor_tags, image, "gcr.io")
    elif repo == 'registry.k8s.io':
        tags_data = get_repo_gcr_tags(harbor_tags, image, "k8s.lixd.xyz")
    elif repo == 'nvcr.io':
        tags_data = get_nvcr_io_tags(harbor_tags, image)
    elif repo == 'quay.io':
        tags_data = get_repo_quay_tags(harbor_tags, image)
    elif repo == "docker.io":
        tags_data = get_docker_io_tags(harbor_tags, image)
    return tags_data


def generate_dynamic_conf():
    """
    生成动态同步配置
    :return:
    """

    print('[generate_dynamic_conf] start.')
    config = None
    with open(CONFIG_FILE, 'r') as stream:
        try:
            config = yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print('[Get Config]', e)
            exit(1)

    print('[config]', config)

    for repo in config['images']:
        skopeo_sync_repo_name = ""
        if repo == "docker.io":
            skopeo_sync_repo_name = "dockerhub.jobcher.com"
        elif repo == "gcr.io":
            skopeo_sync_repo_name = "gcr.io"
        elif repo == "registry.k8s.io":
            skopeo_sync_repo_name = "k8s.lixd.xyz"
        elif repo == "quay.io":
            skopeo_sync_repo_name = "quay.lixd.xyz"
        elif repo == "nvcr.io":
            skopeo_sync_repo_name = "ngc.nju.edu.cn"
        else:
            skopeo_sync_repo_name = repo

        harbor_tags = get_harbor_tags(repo)
        for image in config['images'][repo]:
            skopeo_sync_data = {}
            if repo not in skopeo_sync_data:
                skopeo_sync_data[skopeo_sync_repo_name] = {'images': {}}
            if config['images'][repo] is None:
                continue
            print("[image] {image}".format(image=image))
            sync_tags = get_repo_tags(repo, image, harbor_tags)
            if len(sync_tags) > 0:
                skopeo_sync_data[skopeo_sync_repo_name]['images'][image] = sync_tags
            else:
                print('[{image}] no sync tag.'.format(image=image))

            image_replaced = image.replace('/', '_')
            yaml_name = f'sync_yamls/{repo}_{image_replaced}.yaml'
            with open(yaml_name, 'w+') as f:
                yaml.safe_dump(skopeo_sync_data, f, default_flow_style=False)

        print('[generate_dynamic_conf] done.', end='\n')


def generate_custom_conf():
    """
    生成自定义的同步配置
    :return:
    """

    print('[generate_custom_conf] start.')
    custom_sync_config = None
    with open(CUSTOM_SYNC_FILE, 'r') as stream:
        try:
            custom_sync_config = yaml.safe_load(stream)
        except yaml.YAMLError as e:
            print('[Get Config]', e)
            exit(1)

    print('[custom_sync config]', custom_sync_config)

    custom_skopeo_sync_data = {}

    for repo in custom_sync_config:
        if repo not in custom_skopeo_sync_data:
            custom_skopeo_sync_data[repo] = {'images': {}}
        if custom_sync_config[repo]['images'] is None:
            continue

        harbor_tags = get_harbor_tags(repo)
        for image in custom_sync_config[repo]['images']:
            for tag in custom_sync_config[repo]['images'][image]:
                if tag in harbor_tags:
                    continue
                if image not in custom_skopeo_sync_data[repo]['images']:
                    custom_skopeo_sync_data[repo]['images'][image] = [tag]
                else:
                    custom_skopeo_sync_data[repo]['images'][image].append(tag)

    print('[custom_sync data]', custom_skopeo_sync_data)

    with open(CUSTOM_SYNC_FILE, 'w+') as f:
        yaml.safe_dump(custom_skopeo_sync_data, f, default_flow_style=False)

    print('[generate_custom_conf] done.', end='\n\n')


def generate_sync_image_sh():
    now = datetime.datetime.now()
    now_day = now.strftime("%Y%m%d")
    now_time = now.strftime("%Y%m%d-%H%M%S")
    log_file_path = f'{LOG_FILE_PATH}/{now_day}'
    if not os.path.exists(log_file_path):
        os.mkdir(log_file_path)
    # Read the YAML file
    with open(CONFIG_FILE, 'r') as file:
        images = yaml.safe_load(file)['images']

    for image_source, repos in images.items():
        shell_script_lines = ['skopeo login -u harbor -p GphVjjGe0 harbor.rsq.cn']
        for repo_name in repos:
            src_repo_name = repo_name
            dest_repo_name = '/'.join(src_repo_name.split('/')[:-1]) if '/' in src_repo_name else ""
            image_name = image_source
            src_repo_name_replaced = src_repo_name.replace('/', '_')
            yaml_name = f'{image_name}_{src_repo_name_replaced}.yaml'
            shell_command = f"skopeo --insecure-policy sync -a --keep-going --src yaml --dest docker {yaml_name} harbor.rsq.cn/{image_name}/{dest_repo_name} >> /var/log/sync_images/{now_day}/{now_time}_{image_name}.log 2>&1"
            shell_script_lines.append(shell_command)

        # 保存文件
        shell_script_content = "\n".join(shell_script_lines)
        with open(f'sync_yamls/sync_{image_source}.sh', 'w') as file:
            file.write(shell_script_content)

    print("Shell script has been generated and saved to sync_images.sh")


generate_dynamic_conf()
# generate_custom_conf()
generate_sync_image_sh()
