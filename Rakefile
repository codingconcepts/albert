require 'rake'

orchOutput = 'harambe-orchestrator.exe'
agentOutput = 'harambe-agent.exe'

task :nats do
    Dir.chdir('cmd\\nats') do
        sh("start gnatsd --config nats.conf")
    end
end

task :orch do
    Dir.chdir('cmd\\orchestrator') do
        sh('go', 'build', '-o', orchOutput)
        sh("start #{orchOutput}")
    end
end

task :agent do
    Dir.chdir('cmd\\agent') do
        sh('go', 'build', '-o', agentOutput)
        sh("start #{agentOutput}")
        sh("start #{agentOutput}")
        sh("start #{agentOutput}")
        sh("start #{agentOutput}")
    end
end

task :run do
    Rake::Task["nats"].execute
    Rake::Task["agent"].execute
    Rake::Task["orch"].execute
end

task :stop do
    begin
        sh("taskkill", "/f", "/t", "/im", orchOutput)
    rescue
        puts "error caught stopping #{orchOutput}"
    end

    begin
        sh("taskkill", "/f", "/t", "/im", agentOutput)
    rescue
        puts "error caught stopping #{agentOutput}"
    end

    begin
        sh("taskkill", "/f", "/t", "/im", "gnatsd.exe")
    rescue
        puts "error caught stopping gnatsd.exe"
    end
end
