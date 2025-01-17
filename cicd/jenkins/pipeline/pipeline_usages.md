## 1 Pipeline 批量pull code

```groovy
#!/usr/bin/env groovy
def gitclone(String gitrepo, String gitbranch, String credential) {
    checkout(
        [   
            $class: 'GitSCM', 
            branches: [[name: gitbranch]], 
            doGenerateSubmoduleConfigurations: false, 
            extensions: [
                [$class: 'CheckoutOption', timeout: 60],
                [$class: 'GitLFSPull'],
                [$class: 'CloneOption', noTags: false, reference: '', shallow: false, timeout: 60],
                [$class: 'SubmoduleOption', timeout: 60, disableSubmodules: false, parentCredentials: true, recursiveSubmodules: true, reference: '', trackingSubmodules: false]
            ], 
            submoduleCfg: [], 
            userRemoteConfigs: [[credentialsId: credential, 
            url: gitrepo]]
        ]
    )
}

node(env.BuildMachineLabel){
    stage("Pull Code on win") {
        dir(env.ws_win + 'CloudRendererAgentEdge'){
            gitclone(env.CloudRendererAgentEdge, env.cloud_branch, env.cloud_credential)
        }

        dir(env.ws_win + 'TileService'){
            gitclone(env.TileService, env.master_branch, env.credential_id)
        }

        dir(env.ws_win + 'AES_data_service'){
            gitclone(env.AES_data_service, env.master_branch, env.credential_id)
        }

        dir(env.ws_win + 'std'){
            gitclone(env.std, env.master_branch, env.credential_id)
        }

        dir(env.ws_win + 'AesFramework'){
            gitclone(env.AesFramework, env.master_branch, env.credential_id)
        }

        dir(env.ws_win + 'AES_Common'){
            gitclone(env.AES_Common, env.master_branch, env.credential_id)
        }

    } //  stage

} // node
```

## 2 Pipeline for Kubernetes pod

```groovy
def is_docker_pull = false
if (env.docker_pull == "yes") {is_docker_pull = true}

podTemplate(
	inheritFrom: '51vr-jenkins-slave',
	containers: [
		containerTemplate(
			name: 'xenial-cis', alwaysPullImage: is_docker_pull, workingDir: '/home/jenkins', \
			image: 'images.rsq.local:5000/ubuntu:xenial-cis', ttyEnabled: true, command: 'cat'
		),
		containerTemplate(
			name: '51vr-docker', workingDir: '/home/jenkins', \
			image: 'images.rsq.local:5000/docker:17.03.2-ce', ttyEnabled: true, command: 'cat'
		)
	], // containers
	volumes: [
		hostPathVolume(hostPath: env.cis_linux, mountPath: '/home/jenkins/workspace'),
		hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
		hostPathVolume(hostPath: '/datadisk/cis', mountPath: '/mnt')
	],
	imagePullSecrets: [ 'images-docker-51vr' ]
)

{   
	echo '===========================================================\nChange Set Start\n==========================================================='
	if(env.change_set != '' && env.change_set != null){
	  for (item in env.change_set.split(';')){
		  echo item
	  }
	}
	echo '===========================================================\nChange Set End\n==========================================================='
	
	parallel(
		'win10':{
			try{
				if (env.platform_win10 == 'true'){
					node(env.BuildMachineLabel){
                         stage('Build'){
                            dir(env.ws_win + '\\Tools\\tools'){
                                bat '''
                                    echo "test1"
                                '''
                            }
                            dir(env.ws_win + '\\Tools\\build'){
                                bat '''
                                    echo "test2"
                                '''
                            }
                        }
					} // node - WIN
				} // if
			} // try
			catch (e){
				emailext attachLog: true, 
				body: '''
				<!DOCTYPE html>
				<html>
				<head>
				<meta charset="UTF-8">
				<title>ALPHA - Please ignore - ${ENV, var="JOB_NAME"}- Build Number: ${BUILD_NUMBER}</title>
				</head>
				<body>
					<div>${JELLY_SCRIPT,template="html"}</div>
				</body>
				</html>''', 
				subject: 'Build - ${JOB_NAME} is failed', to: env.mail_list
				throw e
			} //catch
		}, // win10

		'ubuntu16':{
			try{
				if (env.platform_ubuntu16 == 'true'){
					node(POD_LABEL){
                         stage('Build'){
                            container('xenial-cis-realtoedit'){
                                dir('../Tools/RealToEdit/tools'){
                                    sh '''
                                        echo "widnows stage"
                                    '''
                                }
                            }
                        } // stage - Clean Build - ubuntu16
                    }
                }
			} //try
		}, // ubuntu16
        
	) // parallel

} //end
```

## 3 Pipeline 并行parallel

```groovy
stage("Copy stage"){
	parallel{
		stage("stage-1"){
			steps{
				script{
					sh '''
						echo "stage-1"
					'''
				} //script
			} //steps
		} //stage-1

		stage("stage-2"){
			steps{
				script{
					sh '''
						echo "stage-2"
					'''
				} // script
			} //steps

		} //stage-2

		stage("stage-3"){
			steps{
				script{
					if (! fileExists("${work_space}/dockerfiles/core/Core.tar.gz")) {
						echo "ERROR: core tar package is missing!"
					}
					else {
						sh "rm -f ${work_space}/dockerfiles/core/Core.tar.gz"	
					}
					sh '''
						echo "stage-3"
					'''
				} //script
			} //steps
		} //stage-3
	} //parallel
} //copy stage
```

## 4 Pipeline 条件控制

```groovy
//params.copy_to_dst在jenkins job中需要定义为Bool值
stage('Copy To Dst'){
	when{
		//判定条件为真时
		expression{ params.copy_to_dst.toBoolean() }
	}
	steps{
		script{
			sh '''
				echo "copying...."
			'''
		}
	}
}


stage('Copy To Dst'){
	when{
		//判定条件为假时
		expression{ ! params.copy_to_dst.toBoolean() }
	}
	steps{
		script{
			sh '''
				echo "copying...."
			'''
		}
	}
}
```

用if也可以实现条件判断

```groovy
stage('Copy To Dst'){
	when{
		//判定条件为假时
		expression{ ! params.copy_to_dst.toBoolean() }
	}
	steps{
		script{
			dir("$env.ws_win"){
				if (params.copy_to_nas) {
					bat'''
						echo "test"
					'''
				}
			} //dir
		}//script
	}//steps
}
```

## 5 Pipeline switch-case

```groovy
stage('UE Packing'){
	steps{
		script{
			if (env.make_windows_ue == 'true') {
				switch(env.ue_cook){
					case 'cook':
						if(cook_config == 'Shipping') {
							echo "cook"
						}
					break   
					case 'clean_cook':
						// clean up output folder
						bat 'if exist UnrealEngine\\Output rd /s /q UnrealEngine\\Output'

						if(env.ue_pack_config.contains('Development')) {
							echo "clean cook"
						}
					break
				} //switch
			} //if
		} //script 
	} // steps
} //stage

```

## 6 Pipeline for循环

```groovy
stage('UE Packing'){
	steps{
		script{
			for (cook_config in env.ue_pack_config.split(',')) {
				if(cook_config == 'Shipping') {
					echo "Shipping"
				}
			}
		} //script 
	} // steps
} //stage
```



## 注意事项

1 `customWorkspace`如果跟jenkins 配置中的`Pipeline SCM`一块混用，那么流水线后边checkout代码时候要跟上绝对路径，否则会出现代码紊乱的现象。如果用了Pipeline SCM，那么jenkins会先在workspace目录下pull SCM里边定义的git仓库，由于我们用了customWorkspace，那么SCM的仓库就会被pull到customWorkspace目录下边。

```groovy
pipeline {
	agent{
		node{
			label "node01"
			customWorkspace env.work_space 
		    }  
	} // agent
	stages{
		stage("Checkout Codes "){
			steps{
                script{
					//这里的env.BUILD_TOOLS_DIR需要在jenkins job中传参
					dir("$env.BUILD_TOOLS_DIR"){
						checkout(
							[   
								$class: 'GitSCM', 
								branches: [[name: "${params.BUILD_TOOLS_BRANCH}"]], 
								doGenerateSubmoduleConfigurations: false, 
								extensions: [
									[$class: 'CheckoutOption', timeout: 60],
									[$class: 'GitLFSPull'],
									[$class: 'CloneOption', noTags: false, reference: '', shallow: false, timeout: 60],
									[$class: 'SubmoduleOption', timeout: 60, disableSubmodules: false, parentCredentials: true, recursiveSubmodules: true, reference: '', trackingSubmodules: false]
								], 
								submoduleCfg: [], 
								userRemoteConfigs: [[credentialsId: 'tmpdevops', url: 'http://git.rsq.local/rsq/auto-builder.git']]
							]
						)
					}
                } //script
			} // steps
		} // stage

		stage("make"){
			steps{
				script{
					sh '''
						echo nihao
					'''
				}
			} // steps
		} // stage
        
	} // stages
} // pipeline
```

